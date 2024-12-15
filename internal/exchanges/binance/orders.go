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

func (c *Client) CreateTestOrder(symbol, side, orderType, quantity, price string) error {
	endpoint := "/api/v3/order/test"
	queryParams := url.Values{}

	queryParams.Set("symbol", symbol)
	queryParams.Set("side", side)
	queryParams.Set("type", orderType)
	queryParams.Set("quantity", quantity)

	if orderType == "LIMIT" {
		queryParams.Set("price", price)
		queryParams.Set("timeInForce", "GTC")
	}

	queryParams.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	signature := utils.SignRequest(c.SecretKey, queryParams)
	queryParams.Set("signature", signature)

	fullURL := c.BaseURL + endpoint + "?" + queryParams.Encode()

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-MBX-APIKEY", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s\n", resp.StatusCode, string(body))
	}
	fmt.Println("Test order created successfully!")
	return nil
}
