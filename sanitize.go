package main

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var Characters = []string{"á", "é", "í", "ó", "ú", "ü", "ñ"}

func Sanitize(s string) string {
	s, _ = url.QueryUnescape(s)
	s = strings.ToLower(s)
	return s
}

//Escape converts special characters to hexadecimal representations of their Unicode code points. This is the format that DRAE expects rather than simple URL escaped characters.
func Escape(s string) string {
	for _, c := range Characters {
		r, _ := utf8.DecodeRuneInString(c)
		s = strings.Replace(s, c, "%"+strconv.FormatInt(int64(r), 16), -1)
	}
	return s
}
