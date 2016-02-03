package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeWord(word string) []*Entry {
	return Scrape("http://dle.rae.es/srv/search?w="+word, word)
}

func Scrape(url string, word string) []*Entry {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}

	//Look for entries.
	nodes := doc.Find("article")
	//If no entries were found, there is probably a list of links to definitions.
	if nodes.Length() == 0 {
		//Choose the link for the word that is not a verb.
		doc.Find("li").EachWithBreak(func(i int, s *goquery.Selection) bool {
			url, _ = s.Find("a").Attr("href")
			t := strings.TrimSpace(s.Text())
			t = t[:len(t)-1]
			if strings.HasSuffix(t, "r") {
				return true
			}
			if strings.HasSuffix(t, "rse") {
				return true
			}
			return false
		})
		return Scrape("http://dle.rae.es/srv/"+url, word)
	}

	entries := []*Entry{}
	nodes.Each(func(k int, s *goquery.Selection) {
		etymology := s.Find("p.n2").Text()
		defs := []*Definition{}
		vars := []*Variation{}

		s.Find("p[class^='j']").Each(func(i int, s *goquery.Selection) {
			defs = append(defs, ScrapeDefinition(s))
		})

		s.Find("p[class^='k']").Each(func(i int, s *goquery.Selection) {
			vars = append(vars, &Variation{Variation: s.Text()})

			s.NextAll().EachWithBreak(func(_ int, s *goquery.Selection) bool {
				class, _ := s.Attr("class")
				if strings.HasPrefix(class, "l") {
					vars[i].Definitions = append(vars[i].Definitions, &Definition{Definition: s.Text()})
					return false
				}
				if class != "m" {
					return false
				}
				vars[i].Definitions = append(vars[i].Definitions, ScrapeDefinition(s))
				return true
			})
		})

		entry := &Entry{
			Word:        word,
			Etymology:   etymology,
			Definitions: defs,
			Variations:  vars,
		}

		entries = append(entries, entry)
	})

	return entries
}

func ScrapeDefinition(s *goquery.Selection) *Definition {
	category, _ := s.Find("abbr.g").First().Attr("title")

	return &Definition{
		Category:   category,
		Definition: JoinNodesWithSpace(s.Children().First().NextAll().Not("abbr").Not("span.h")),
		Origin:     ScrapeOrigins(s),
		Notes:      ScrapeNotes(s),
		Examples:   ScrapeExamples(s),
	}
}

func ScrapeOrigins(s *goquery.Selection) []string {
	origins := []string{}
	s.Find("abbr.c").Each(func(i int, s *goquery.Selection) {
		origin, _ := s.Attr("title")
		origins = append(origins, origin)
	})
	return origins
}

func ScrapeNotes(s *goquery.Selection) []string {
	notes := []string{}
	s.Find("abbr.d").Each(func(i int, s *goquery.Selection) {
		note, _ := s.Attr("title")
		notes = append(notes, note)
	})
	return notes
}

func ScrapeExamples(s *goquery.Selection) []string {
	examples := []string{}
	s.Find("span.h").Each(func(i int, s *goquery.Selection) {
		examples = append(examples, s.Text())
	})
	return examples
}

func JoinNodesWithSpace(s *goquery.Selection) string {
	texts := []string{}
	s.Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	return strings.Join(texts, " ")
}
