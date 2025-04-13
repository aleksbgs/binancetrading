package service

import (
	"binancetrading/internal/domain"
	"binancetrading/internal/interfaces"
	"log"
)

type CandlestickService struct {
	aggregators     map[string]*domain.CandlestickAggregator
	exchange        interfaces.Exchange
	storage         interfaces.Storage
	candlestickChan chan *domain.OHLC
}

func NewCandlestickService(exchange interfaces.Exchange, storage interfaces.Storage) *CandlestickService {
	return &CandlestickService{
		aggregators:     make(map[string]*domain.CandlestickAggregator),
		exchange:        exchange,
		storage:         storage,
		candlestickChan: make(chan *domain.OHLC, 100),
	}
}

func (s *CandlestickService) StartAggregation(symbols []string) error {
	tradeChan := make(chan domain.Trade)

	for _, symbol := range symbols {
		aggregator := domain.NewCandlestickAggregator(symbol, s.saveAndBroadcastCandlestick)
		s.aggregators[symbol] = aggregator
	}

	if err := s.exchange.SubscribeToTrades(symbols, tradeChan); err != nil {
		return err
	}

	go func() {
		for trade := range tradeChan {
			if aggregator, exists := s.aggregators[trade.Symbol]; exists {
				if err := aggregator.Update(trade); err != nil {
					log.Printf("Failed to update candlestick for %s: %v", trade.Symbol, err)
				}
			}
		}
	}()

	return nil
}

func (s *CandlestickService) saveAndBroadcastCandlestick(symbol string, candle *domain.OHLC) error {
	if err := s.storage.SaveCandlestick(symbol, candle); err != nil {
		log.Printf("Failed to save candlestick for %s: %v", symbol, err)
		return err
	}

	select {
	case s.candlestickChan <- candle:
	default:
		log.Printf("Candlestick channel full, dropping candlestick for %s", symbol)
	}

	return nil
}

func (s *CandlestickService) CandlestickChan() <-chan *domain.OHLC {
	return s.candlestickChan
}

func (s *CandlestickService) Wait() {
	s.exchange.Wait()
}
