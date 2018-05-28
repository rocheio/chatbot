package main

import (
	"fmt"
	"strings"
)

// Lexicon defines a vocabulary for text structures
type Lexicon struct {
	twoWordFollowers map[string][]string
}

// NewLexicon returns a Lexicon ready to ingest data
func NewLexicon() Lexicon {
	return Lexicon{
		twoWordFollowers: make(map[string][]string),
	}
}

// IngestString adds a strings contents to this Lexicon
func (l Lexicon) IngestString(s string) {
	parts := strings.Split(s, " ")
	if len(parts) < 2 {
		return
	}
	for i := 0; i < len(parts)-2; i++ {
		key := strings.Join([]string{parts[i], parts[i+1]}, " ")
		l.twoWordFollowers[key] = append(l.twoWordFollowers[key], parts[i+2])
	}
}

// SuggestTwoWordFollower returns a following word for two words
func (l Lexicon) SuggestTwoWordFollower(s string) (string, error) {
	val, ok := l.twoWordFollowers[s]
	if !ok || len(val) == 0 {
		return "", fmt.Errorf("two word follower not found for %s", s)
	}
	return val[0], nil
}

func main() {}
