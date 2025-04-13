package domain

import "time"

// Trade represents a trade event from an exchange
type Trade struct {
	Symbol    string
	Price     string
	TradeTime time.Time
}
