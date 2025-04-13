package exchange

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"binancetrading/internal/domain"
	"github.com/adshao/go-binance/v2"
)

// BinanceClient implements the Exchange interface for Binance
type BinanceClient struct {
	wg         sync.WaitGroup
	retryDelay time.Duration
	maxRetries int
}

func NewBinanceClient(retryDelay time.Duration, maxRetries int) *BinanceClient {
	return &BinanceClient{
		retryDelay: retryDelay,
		maxRetries: maxRetries,
	}
}

// SubscribeToTrades subscribes to aggregate trade streams for the given symbols
func (c *BinanceClient) SubscribeToTrades(symbols []string, tradeChan chan<- domain.Trade) error {
	for _, symbol := range symbols {
		c.wg.Add(1)
		go func(sym string) {
			defer c.wg.Done()
			c.connectWithRetry(sym, tradeChan)
		}(symbol)
	}
	return nil
}

// connectWithRetry attempts to connect to the WebSocket stream with retry logic
func (c *BinanceClient) connectWithRetry(symbol string, tradeChan chan<- domain.Trade) {
	attempt := 0
	for {
		if attempt >= c.maxRetries {
			log.Printf("Max retries reached for %s, giving up", symbol)
			return
		}

		doneC, stopC, err := c.connect(symbol, tradeChan)
		if err != nil {
			log.Printf("Failed to connect to %s WebSocket: %v", symbol, err)
			attempt++
			delay := c.retryDelay * time.Duration(1<<uint(attempt))
			jitter := time.Duration(rand.Int63n(int64(delay))) / 2
			log.Printf("Retrying in %v (attempt %d/%d)", delay+jitter, attempt, c.maxRetries)
			time.Sleep(delay + jitter)
			continue
		}

		attempt = 0
		select {
		case <-doneC:
			log.Printf("WebSocket for %s closed, attempting to reconnect", symbol)
			attempt++
			continue
		case <-stopC:
			log.Printf("WebSocket for %s stopped", symbol)
			return
		}
	}
}

// connect establishes a WebSocket connection for a single symbol using the aggTrade stream
func (c *BinanceClient) connect(symbol string, tradeChan chan<- domain.Trade) (doneC, stopC chan struct{}, err error) {
	doneC = make(chan struct{})
	stopC = make(chan struct{})

	tradeHandler := func(event *binance.WsAggTradeEvent) {
		trade := domain.Trade{
			Symbol:    event.Symbol,
			Price:     event.Price,
			TradeTime: time.UnixMilli(event.TradeTime),
		}
		tradeChan <- trade
	}

	errHandler := func(err error) {
		log.Printf("Binance WebSocket error for %s: %v", symbol, err)
	}

	_, _, err = binance.WsAggTradeServe(symbol, tradeHandler, errHandler)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		<-stopC
		close(doneC)
	}()

	return doneC, stopC, nil
}

// Wait waits for all WebSocket connections to close
func (c *BinanceClient) Wait() {
	c.wg.Wait()
}
