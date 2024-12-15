package models

type TickerPrice struct {
	Symbol string `json:"symbol"` // Name of the trading pair
	Price  string `json:"price"`  // Cost per unit
}

type OHLC struct {
	Symbol    string  // Name of the trading pair
	Open      float64 // Initial price (Open)
	High      float64 // Highest price
	Low       float64 // Lowest price
	Close     float64 // Final price (Close)
	Timestamp int64   // Time
}
