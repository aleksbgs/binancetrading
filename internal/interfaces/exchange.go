package interfaces

import "binancetrading/internal/domain"

// Exchange defines the interface for interacting with an exchange
type Exchange interface {
	SubscribeToTrades(symbols []string, tradeChan chan<- domain.Trade) error
	Wait()
}
