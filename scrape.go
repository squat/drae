package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Scrape(word string) *Entry {
	r1, err := http.Get("http://lema.rae.es/drae/srv/search?val=" + word)

	if err != nil {
		panic(err)
	}

	r2 := Solve(r1)
	defer r2.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(r2)
	if err != nil {
		panic(err)
	}

	node := doc.Find("div").First()
	delimiter := node.Children().Filter("p:not([class])").First()

	etymology := strings.TrimSpace(node.Find("span.a").Text())
	defs := []*Definition{}

	delimiter.NextAll().EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.HasClass("p") {
			return false
		}
		category, _ := s.Find("span[title]").First().Attr("title")
		origins := []string{}
		notes := []string{}
		examples := []string{}

		s.Find("span.d i span.d[title]").Each(func(i int, s *goquery.Selection) {
			origin, _ := s.Attr("title")
			origins = append(origins, origin)
		})

		s.Clone().Find("span[title]").First().Remove().End().End().Find("span.d i span.d[title]").Remove().End().Find("span.d[title]").Each(func(i int, s *goquery.Selection) {
			note, _ := s.Attr("title")
			notes = append(notes, note)
		})

		s.Find("span.h i").Each(func(i int, s *goquery.Selection) {
			examples = append(examples, s.Text())
		})

		def := &Definition{
			Category:   category,
			Definition: strings.TrimSpace(s.Find("span.b").Clone().Children().Remove().End().Text()),
			Origin:     origins,
			Notes:      notes,
			Examples:   examples,
		}

		defs = append(defs, def)
		return true
	})

	entry := &Entry{
		Word:        word,
		Etymology:   etymology,
		Definitions: defs,
	}

	return entry
}
