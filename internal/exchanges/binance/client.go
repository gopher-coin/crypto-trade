package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gopher-coin/crypto-trade/internal/utils"
	"github.com/gopher-coin/crypto-trade/pkg/models"
)

type Client struct {
	APIKey    string
	SecretKey string
	BaseURL   string
}

type Kline struct {
	OpenTime  int64
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	CloseTime int64
}

func NewClient(apiKey, secretKey, baseURL string) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   baseURL,
	}
}

func (c *Client) GetAccountInfo() ([]models.Balance, error) {
	endpoint := "/api/v3/account"
	queryParams := url.Values{}

	queryParams.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	signature := utils.SignRequest(c.SecretKey, queryParams)
	queryParams.Set("signature", signature)

	fullURL := c.BaseURL + endpoint + "?" + queryParams.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w\n", err)
	}

	req.Header.Set("X-MBX-APIKEY", c.APIKey)

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

	var accountResponse struct {
		Balances []models.Balance `json:"balances"`
	}
	if err := json.Unmarshal(body, &accountResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return accountResponse.Balances, nil
}

func (c *Client) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	endpoint := "/api/v3/klines"

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)
	params.Set("limit", fmt.Sprintf("%d", limit))

	fullURL := c.BaseURL + endpoint + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch klines: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %s\n", string(body))
	}

	var rawData [][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decoderesponse: %w\n", err)
	}

	var klines []Kline
	for _, k := range rawData {
		klines = append(klines, Kline{
			OpenTime:  int64(k[0].(float64)),
			Open:      k[1].(string),
			High:      k[2].(string),
			Low:       k[3].(string),
			Close:     k[4].(string),
			Volume:    k[5].(string),
			CloseTime: int64(k[6].(float64)),
		})
	}
	return klines, nil
}
