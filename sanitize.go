package main

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const characters string = "áéíóúüñ"

func Sanitize(s string) string {
	s, _ = url.QueryUnescape(s)
	s = strings.ToLower(s)
	return s
}

//Escape converts special characters to hexadecimal representations of their Unicode code points. This is the format that DRAE expects rather than simple URL escaped characters.
func Escape(s string) string {
	for i, w := 0, 0; i < len(characters); i += w {
		r, width := utf8.DecodeRuneInString(characters[i:])
		c := characters[i : i+width]
		s = strings.Replace(s, c, "%"+strconv.FormatInt(int64(r), 16), -1)
		w = width
	}
	return s
}
