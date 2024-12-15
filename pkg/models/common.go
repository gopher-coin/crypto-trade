package models

type ErrorResponse struct {
	Code int    `json:"code"` // API error code
	Msg  string `json:"msg"`  // API error message
}

type Balance struct {
	Asset  string `json:"asset"`  // Asset name
	Free   string `json:"free"`   // Available quantity
	Locked string `json:"locked"` // Blocked quantity
}
