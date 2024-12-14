package errors

import "fmt"

func HandleAPIError(code int, msg string) error {
	return fmt.Errorf("API error %d: %s \n", code, msg)
}
