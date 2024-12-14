package binance

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gopher-coin/crypto-trade/internal/utils"
)

type Client struct {
	APIKey    string
	SecretKey string
	BaseURL   string
}

func NewClient(apiKey, secretKey, baseURL string) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   baseURL,
	}
}

func (c *Client) GetAccountInfo() ([]byte, error) {
	endpoint := "/api/v3/openOrders"
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

	return body, nil
}
