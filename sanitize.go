package main

import (
	"net/url"
	"strings"
)

func Sanitize(s string) string {
	s, _ = url.QueryUnescape(s)
	s = strings.ToLower(s)
	return s
}
