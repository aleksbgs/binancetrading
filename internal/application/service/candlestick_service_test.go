package service

import (
	"testing"
	"time"

	"binancetrading/internal/domain"
	"binancetrading/internal/interfaces"
)

func TestCandlestickService_StartAggregation_Simple(t *testing.T) {
	// Initialize mock dependencies
	mockExchange := interfaces.NewMockExchange()
	mockStorage := interfaces.NewMockStorage()

	// Initialize service
	service := NewCandlestickService(mockExchange, mockStorage)

	// Start aggregation for a symbol
	symbols := []string{"BTCUSDT"}
	if err := service.StartAggregation(symbols); err != nil {
		t.Fatalf("StartAggregation failed: %v", err)
	}

	// Send a trade through the mock exchange
	startTime := time.Date(2025, 4, 12, 19, 4, 0, 0, time.UTC)
	trade := domain.Trade{
		Symbol:    "BTCUSDT",
		Price:     "65000.0",
		TradeTime: startTime,
	}
	mockExchange.SendTrade(trade)

	// Send another trade in a new minute to trigger saving
	trade2 := domain.Trade{
		Symbol:    "BTCUSDT",
		Price:     "65100.0",
		TradeTime: startTime.Add(1 * time.Minute),
	}
	mockExchange.SendTrade(trade2)

	// Allow some time for processing
	time.Sleep(100 * time.Millisecond)

	// Verify that the candlestick was saved
	savedCandles := mockStorage.GetSavedCandlesticks()
	if len(savedCandles) != 1 {
		t.Fatalf("Expected 1 candlestick to be saved, got %d", len(savedCandles))
	}
	saved := savedCandles[0]
	if saved.Symbol != "BTCUSDT" { // Use exported field Symbol
		t.Errorf("Expected symbol to be BTCUSDT, got %s", saved.Symbol)
	}
	if saved.Candle.Open != "65000.0" || saved.Candle.Close != "65000.0" { // Use exported field Candle
		t.Errorf("Expected saved Open and Close to be 65000.0, got Open: %s, Close: %s", saved.Candle.Open, saved.Candle.Close)
	}

	// Verify that the candlestick was broadcasted
	select {
	case candle := <-service.CandlestickChan():
		if candle.Open != "65000.0" || candle.Close != "65000.0" {
			t.Errorf("Expected broadcasted Open and Close to be 65000.0, got Open: %s, Close: %s", candle.Open, candle.Close)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected a candlestick to be broadcasted, but none received")
	}
}
