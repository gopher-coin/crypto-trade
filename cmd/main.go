package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey    string
	SecretKey string
	BaseURL   string
}

type TickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type errorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	env := os.Getenv("BINANCE_ENV")

	if env == "LIVE" {
		return &Config{
			APIKey:    os.Getenv("LIVE_BINANCE_API_KEY"),
			SecretKey: os.Getenv("LIVE_BINANCE_SECRET_KEY"),
			BaseURL:   os.Getenv("LIVE_BINANCE_BASE_URL"),
		}, nil
	}

	return &Config{
		APIKey:    os.Getenv("TEST_BINANCE_API_KEY"),
		SecretKey: os.Getenv("TEST_BINANCE_SECRET_KEY"),
		BaseURL:   os.Getenv("TEST_BINANCE_BASE_URL"),
	}, nil
}

func signRequest(secretKey string, queryParams url.Values) string {

	queryString := queryParams.Encode()

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(queryString))

	return hex.EncodeToString(h.Sum(nil))
}

func getAccountInfo(config *Config) ([]byte, error) {
	endpoint := "/api/v3/openOrders"
	queryParams := url.Values{}

	queryParams.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	signature := signRequest(config.SecretKey, queryParams)
	queryParams.Set("signature", signature)

	fullURL := config.BaseURL + endpoint + "?" + queryParams.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}

	req.Header.Set("X-MBX-APIKEY", config.APIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w\n", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s\n", resp.StatusCode, string(body))
	}
	return body, nil
}

func BuildPriceMap(prices []TickerPrice) (map[string]string, error) {
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
		var apiError errorResponse
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

	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	fmt.Printf("Connected to %s environment\n", os.Getenv("BINANCE_ENV"))

	body, err := getAccountInfo(config)
	if err != nil {
		fmt.Println("Error getting account info:\n", err)
		return
	}

	fmt.Println("Account Info:\n", string(body))

	//url := "https://api.binance.com/api/v3/ticker/price"
	//url := "https://httpbin.org/status/500"

	// body, err := FetchBinanceDataWithRetry(url, 3, 2*time.Second)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// fmt.Println("Response body:", string(body[:200]))

	// var prices []TickerPrice
	// if err := json.Unmarshal(body, &prices); err != nil {
	// 	fmt.Printf("Failed to decode prices: %v\n", err)
	// 	return
	// }

	// priceMap, err := BuildPriceMap(prices)
	// if err != nil {
	// 	fmt.Println("error buiding price map:", err)
	// 	return
	// }

	// symbol := "BTCUSDT"
	// price, err := GetTickerPrice(priceMap, symbol)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Printf("Price for %s: %s\n", symbol, price)
	// }
}
