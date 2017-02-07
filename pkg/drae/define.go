package drae

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const raeAPI = "http://dle.rae.es/srv/"

// NotFoundError is a special error type that is returned when the word
// could not be found in the API.
type NotFoundError struct {
	err error
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("failed to find word: %v", e.err)
}

// Define is a wrapper for Scrape that takes a word you'd like to define and
// returns a slice of Entry's.
func Define(word string) ([]*Entry, error) {
	word, err := Sanitize(word)
	if err != nil {
		return nil, err
	}
	return scrape(raeAPI+"search?w="+word, word)
}

// scrape takes a URL corresponding to a resource on RAE and returns a slice of
// Entry's.
func scrape(url string, word string) ([]*Entry, error) {
	if word == "" {
		return nil, errors.New("word cannot be empty")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate request to rae: %v", err)
	}
	req.Header.Set("User-Agent", "")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to rae: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	solveRe := regexp.MustCompile("Please enable JavaScript to view the page content.")
	if solveRe.Match(body) {
		// We need to do some work to get an actual response.
		body, err = solve(body)
		if err != nil {
			return nil, fmt.Errorf("failed to solve challenge: %v", err)
		}
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response from rae: %v", err)
	}

	// Look for entries.
	nodes := doc.Find("article")
	// If no entries were found, there is probably a list of links to definitions.
	if nodes.Length() == 0 {
		// Choose the link for the word that is not a verb.
		url, err := findNonVerbURL(doc.Find("body"))
		if err != nil {
			return nil, NotFoundError{err}
		}
		return scrape(raeAPI+url, word)
	}

	var wg sync.WaitGroup
	wg.Add(nodes.Length())
	entries := make([]*Entry, nodes.Length())
	nodes.Each(func(k int, s *goquery.Selection) {
		// Parallelize generating entries.
		go func(k int, s *goquery.Selection) {
			etymology := s.Find("p.n2").Text()
			var defs []*Definition
			var vars []*Variation

			s.Find("p[class^='j']").Each(func(i int, s *goquery.Selection) {
				defs = append(defs, scrapeDefinition(s))
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
					vars[i].Definitions = append(vars[i].Definitions, scrapeDefinition(s))
					return true
				})
			})

			// Insert directly into entries vs using a channel to guarantee order.
			entries[k] = &Entry{
				Word:        word,
				Etymology:   etymology,
				Definitions: defs,
				Variations:  vars,
			}
			wg.Done()
		}(k, s)
	})

	wg.Wait()
	return entries, nil
}

func findNonVerbURL(s *goquery.Selection) (string, error) {
	var url string
	s.Find("li").EachWithBreak(func(i int, s *goquery.Selection) bool {
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

	if url == "" {
		return "", errors.New("failed to find url")
	}

	return url, nil
}
func scrapeDefinition(s *goquery.Selection) *Definition {
	category, _ := s.Find("abbr").First().Attr("title")

	return &Definition{
		Category:   category,
		Definition: joinNodesWithSpace(s.Children().First().NextAll().Not("abbr").Not("span.h")),
		Origin:     scrapeOrigins(s),
		Notes:      scrapeNotes(s),
		Examples:   scrapeExamples(s),
	}
}

func scrapeOrigins(s *goquery.Selection) []string {
	var origins []string
	s.Find("abbr.c").Each(func(i int, s *goquery.Selection) {
		origin, _ := s.Attr("title")
		origins = append(origins, origin)
	})
	return origins
}

func scrapeNotes(s *goquery.Selection) []string {
	var notes []string
	s.Find("abbr").Not("abbr:first-of-type").Not("abbr.c").Each(func(i int, s *goquery.Selection) {
		note, _ := s.Attr("title")
		notes = append(notes, note)
	})
	return notes
}

func scrapeExamples(s *goquery.Selection) []string {
	var examples []string
	s.Find("span.h").Each(func(i int, s *goquery.Selection) {
		examples = append(examples, s.Text())
	})
	return examples
}

func joinNodesWithSpace(s *goquery.Selection) string {
	var texts []string
	s.Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	return strings.Join(texts, " ")
}
