package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetTickerPrice(prices []TickerPrice, value string) (string, error) {
	value = strings.ToUpper(value)
	for _, ticker := range prices {
		if ticker.Symbol == value {
			return ticker.Price, nil
		}
	}
	return "", fmt.Errorf("symbol %s not found", value)
}

func main() {
	url := "https://api.binance.com/api/v3/ticker/price"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var prices []TickerPrice
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		fmt.Printf("Error decoding response data: %v\n", err)
		return
	}

	for i, ticker := range prices {
		if i >= 5 {
			break
		}
		fmt.Printf("%s: %s\n", ticker.Symbol, ticker.Price)
	}
	btcUsdt, err := GetTickerPrice(prices, "BTCUSDT")
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	fmt.Println(btcUsdt)
}
