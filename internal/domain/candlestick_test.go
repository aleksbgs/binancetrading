package domain

import (
	"testing"
	"time"
)

// TestCandlestickAggregator_Update_Simple tests the basic functionality of the Update method
func TestCandlestickAggregator_Update_Simple(t *testing.T) {
	// Mock save function to capture saved candlesticks
	var savedCandles []OHLC
	saveFunc := func(symbol string, candle *OHLC) error {
		savedCandles = append(savedCandles, *candle)
		return nil
	}

	// Initialize aggregator
	aggregator := NewCandlestickAggregator("BTCUSDT", saveFunc)

	// Test case 1: First trade in a minute
	startTime := time.Date(2025, 4, 12, 19, 4, 0, 0, time.UTC)
	trade1 := Trade{
		Symbol:    "BTCUSDT",
		Price:     "65000.0",
		TradeTime: startTime,
	}
	if err := aggregator.Update(trade1); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the current candlestick
	if aggregator.current == nil {
		t.Fatal("Expected current candlestick to be initialized")
	}
	if aggregator.current.Open != "65000.0" || aggregator.current.Close != "65000.0" {
		t.Errorf("Expected Open and Close to be 65000.0, got Open: %s, Close: %s", aggregator.current.Open, aggregator.current.Close)
	}

	// Test case 2: New minute, should save the previous candlestick
	trade2 := Trade{
		Symbol:    "BTCUSDT",
		Price:     "65100.0",
		TradeTime: startTime.Add(1 * time.Minute),
	}
	if err := aggregator.Update(trade2); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the saved candlestick
	if len(savedCandles) != 1 {
		t.Fatalf("Expected 1 candlestick to be saved, got %d", len(savedCandles))
	}
	saved := savedCandles[0]
	if saved.Open != "65000.0" || saved.Close != "65000.0" {
		t.Errorf("Expected saved Open and Close to be 65000.0, got Open: %s, Close: %s", saved.Open, saved.Close)
	}

	// Verify the new current candlestick
	if aggregator.current.Open != "65100.0" || aggregator.current.Close != "65100.0" {
		t.Errorf("Expected new Open and Close to be 65100.0, got Open: %s, Close: %s", aggregator.current.Open, aggregator.current.Close)
	}
}
