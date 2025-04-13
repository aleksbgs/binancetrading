package domain

import (
	"sync"
	"time"
)

// OHLC represents a 1-minute candlestick
type OHLC struct {
	Symbol    string
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	Timestamp time.Time
}

// CandlestickAggregator aggregates trades into OHLC candlesticks
type CandlestickAggregator struct {
	symbol    string
	current   *OHLC
	trades    []string
	startTime time.Time
	mu        sync.Mutex
	saveFunc  func(symbol string, candle *OHLC) error // Callback to save candlestick
}

func NewCandlestickAggregator(symbol string, saveFunc func(symbol string, candle *OHLC) error) *CandlestickAggregator {
	return &CandlestickAggregator{
		symbol:   symbol,
		saveFunc: saveFunc,
	}
}

// Update processes a new trade and updates the OHLC
func (ca *CandlestickAggregator) Update(trade Trade) error {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	// Truncate to the start of the current minute
	currentMinute := trade.TradeTime.Truncate(time.Minute)

	if ca.current == nil || currentMinute.After(ca.startTime) {
		// Finalize the previous candlestick
		if ca.current != nil {
			if err := ca.saveFunc(ca.symbol, ca.current); err != nil {
				return err
			}
		}

		// Start a new candlestick
		ca.current = &OHLC{
			Symbol:    ca.symbol,
			Open:      trade.Price,
			High:      trade.Price,
			Low:       trade.Price,
			Close:     trade.Price,
			Timestamp: currentMinute,
			Volume:    trade.Price,
		}
		ca.startTime = currentMinute
		ca.trades = []string{trade.Price}
	} else {
		// Update the current candlestick
		ca.trades = append(ca.trades, trade.Price)
		if trade.Price > ca.current.High {
			ca.current.High = trade.Price
		}
		if trade.Price < ca.current.Low {
			ca.current.Low = trade.Price
		}
		ca.current.Close = trade.Price
	}
	return nil
}
