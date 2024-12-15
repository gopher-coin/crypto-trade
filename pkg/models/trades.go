package models

type Trade struct {
	ID        int64   // Transaction ID
	Symbol    string  // Name of the trading pair
	Side      string  // BUY or SELL
	Quantity  float64 // Number of units
	Price     float64 // Cost per unit
	Timestamp int64   // Transaction time
}
