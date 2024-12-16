package main

import (
	"context"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gopher-coin/crypto-trade/internal/config"
	"github.com/gopher-coin/crypto-trade/internal/db"
	"github.com/gopher-coin/crypto-trade/internal/db/queries"
	"github.com/gopher-coin/crypto-trade/internal/exchanges/binance"
	"github.com/gopher-coin/crypto-trade/pkg/models"
)

func BuildPriceMap(prices []models.TickerPrice) (map[string]string, error) {
	if len(prices) == 0 {
		return nil, fmt.Errorf("input price list is empty")
	}
	pricesMap := make(map[string]string)
	for _, ticker := range prices {
		pricesMap[ticker.Symbol] = ticker.Price
	}
	return pricesMap, nil
}

func GetTickerPrice(pricesMap map[string]string, value string) (string, error) {
	price, ok := pricesMap[strings.ToUpper(value)]
	if !ok {
		return "", fmt.Errorf("symbol %s not found", value)
	}
	return price, nil
}

func FetchBinanceDataWithErrors(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		var apiError models.ErrorResponse
		if err := json.Unmarshal(body, &apiError); err == nil && apiError.Msg != "" {
			return nil, fmt.Errorf("API error %d: %s", apiError.Code, apiError.Msg)
		}
		return nil, fmt.Errorf("HTTP error %d: %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w\n", err)
	}

	return body, nil
}

func FetchBinanceDataWithRetry(url string, maxRetries int, retryDelay time.Duration) ([]byte, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		body, err := FetchBinanceDataWithErrors(url)
		if err == nil {
			return body, nil
		}
		if strings.Contains(err.Error(), "HTTP error 5") {
			lastErr = err
			fmt.Printf("Retry %d/%d: %s\n", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}
		return nil, err
	}
	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v\n", err)
	}
	defer db.CloseDB()

	client := binance.NewClient(cfg.APIKey, cfg.SecretKey, cfg.BaseURL)

	intervals := []string{"1m", "5m", "15m", "1h", "1d"}
	symbol := "ETHUSDT"

	for _, interval := range intervals {
		klines, err := client.GetKlines(symbol, interval, 10)
		if err != nil {
			log.Printf("Failed to fetch klines for %s %s: %v\n", symbol, interval, err)
			continue
		}
		ohlcQueries := &queries.OHLCQueries{Pool: db.Pool}

		for _, kline := range klines {
			ohlc := models.OHLC{
				Symbol:    symbol,
				Open:      kline.Open,
				High:      kline.High,
				Low:       kline.Low,
				Close:     kline.Close,
				Volume:    kline.Volume,
				Timestamp: kline.OpenTime,
			}
			err := ohlcQueries.InsertOHLC(context.Background(), ohlc)
			if err != nil {
				log.Printf("Failed to insert OHLC data: %v", err)
			}
		}
		fmt.Printf("Inserted OHLC data for interval: %s\n", interval)
	}
	fmt.Println("OHLC data fetching and saving complete.")

	// err := db.ConnectDB()
	// if err != nil {
	// 	log.Fatalf("Database connection failed: %v\n", err)
	// }
	// defer db.CloseDB()
	// log.Println("Aplication is running...")

	// client := binance.NewClient(cfg.APIKey, cfg.SecretKey, cfg.BaseURL)

	// balances, err := client.GetAccountInfo()
	// if err != nil {
	// 	fmt.Println("Error getting account info:", err)
	// 	return
	// }

	// for _, balance := range balances {
	// 	fmt.Printf("Asset: %s, Free: %s, Locked: %s\n", balance.Asset, balance.Free, balance.Locked)
	// }

	// err = client.CreateTestOrder("BTCUSDT", "BUY", "LIMIT", "0.001", "90000")
	// if err != nil {
	// 	fmt.Println("Error creating test order:", err)
	// 	return
	// }

	// fmt.Println("Test order executed successfully!")

}
