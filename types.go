package main

type Definition struct {
	Category   string   `json:"category"`
	Definition string   `json:"definition"`
	Origin     []string `json:"origin"`
	Notes      []string `json:"notes"`
	Examples   []string `json:"examples"`
}

type Entry struct {
	Word        string        `json:"word"`
	Etymology   string        `json:"etylmology"`
	Definitions []*Definition `json:"definitions"`
}
