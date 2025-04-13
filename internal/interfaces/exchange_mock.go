package interfaces

import (
	"sync"

	"binancetrading/internal/domain"
)

// MockExchange is a mock implementation of the Exchange interface
type MockExchange struct {
	wg        sync.WaitGroup
	tradeChan chan<- domain.Trade // Change to send-only channel
	symbols   []string
}

func NewMockExchange() *MockExchange {
	return &MockExchange{}
}

func (m *MockExchange) SubscribeToTrades(symbols []string, tradeChan chan<- domain.Trade) error {
	m.symbols = symbols
	m.tradeChan = tradeChan
	return nil
}

func (m *MockExchange) Wait() {
	m.wg.Wait()
}

// SendTrade simulates sending a trade to the trade channel
func (m *MockExchange) SendTrade(trade domain.Trade) {
	m.tradeChan <- trade // Sending to a send-only channel is allowed
}
