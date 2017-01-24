package drae

import (
	"fmt"
	"net/url"
	"strings"
)

// Sanitize cleans a string so that it will be accepted by the rae website.
func Sanitize(s string) (string, error) {
	s, err := url.QueryUnescape(s)
	if err != nil {
		return "", fmt.Errorf("failed to sanitize the provided word: %v", err)
	}
	s = strings.ToLower(s)
	return s, nil
}
