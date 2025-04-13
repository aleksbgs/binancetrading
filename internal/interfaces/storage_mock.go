package interfaces

import (
	"binancetrading/internal/domain"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	candlesticks []struct {
		Symbol string       // Export by capitalizing
		Candle *domain.OHLC // Export by capitalizing
	}
	getResults []domain.OHLC
	getError   error
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) SaveCandlestick(symbol string, candle *domain.OHLC) error {
	m.candlesticks = append(m.candlesticks, struct {
		Symbol string
		Candle *domain.OHLC
	}{Symbol: symbol, Candle: candle})
	return nil
}

func (m *MockStorage) GetCandlesticks(symbol, startTime, endTime string) ([]domain.OHLC, error) {
	return m.getResults, m.getError
}

func (m *MockStorage) Close() error {
	return nil
}

// GetSavedCandlesticks returns the saved candlesticks for inspection
func (m *MockStorage) GetSavedCandlesticks() []struct {
	Symbol string
	Candle *domain.OHLC
} {
	return m.candlesticks
}
