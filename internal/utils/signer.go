package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
)

func SignRequest(secretKey string, queryParams url.Values) string {

	queryString := queryParams.Encode()

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(queryString))

	return hex.EncodeToString(h.Sum(nil))
}
