package interfaces

import "binancetrading/internal/domain"

// Storage defines the interface for persisting candlesticks
type Storage interface {
	SaveCandlestick(symbol string, candle *domain.OHLC) error
	GetCandlesticks(symbol, startTime, endTime string) ([]domain.OHLC, error)
	Close() error
}
