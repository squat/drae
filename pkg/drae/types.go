package drae

// Definition represents a single definition for a word, including the
// category, origin etc.
type Definition struct {
	Category   string   `json:"category"`
	Definition string   `json:"definition"`
	Origin     []string `json:"origin"`
	Notes      []string `json:"notes"`
	Examples   []string `json:"examples"`
}

// Variation represents an alternative way to use the word, e.g. in a different context.
type Variation struct {
	Definitions []*Definition `json:"definitions"`
	Variation   string        `json:"variation"`
}

// Entry represents a collection of definitions and variations for a given word.
type Entry struct {
	Word        string        `json:"word"`
	Etymology   string        `json:"etymology"`
	Definitions []*Definition `json:"definitions"`
	Variations  []*Variation  `json:"variations"`
}
