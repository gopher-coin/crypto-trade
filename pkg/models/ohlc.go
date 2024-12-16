package models

type TickerPrice struct {
	Symbol string `json:"symbol"` // Name of the trading pair
	Price  string `json:"price"`  // Cost per unit
}

type OHLC struct {
	Symbol    string // Name of the trading pair
	Open      string // Initial price (Open)
	High      string // Highest price
	Low       string // Lowest price
	Close     string // Final price (Close)
	Volume    string // Trading volume for a certain period
	Timestamp int64  // Time
}
