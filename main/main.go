package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// normalize returns a string in lowercase with only spaces for whitespace
func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// alphanum returns a string with only alphanumeric characters
func alphanum(s string) string {
	var clean string
	for _, r := range s {
		if unicode.IsLetter(r) {
			clean += string(r)
		}
	}
	return clean
}

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
		val := alphanum(parts[i+2])
		l.twoWordFollowers[key] = append(l.twoWordFollowers[key], val)
	}
}

// IngestReader adds all text from a Reader to this Lexicon
func (l Lexicon) IngestReader(r io.Reader) {
	buf := bufio.NewReader(r)
	for {
		sentence, err := buf.ReadString([]byte(".")[0])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		l.IngestString(normalize(sentence))
	}
}

// SuggestTwoWordFollower returns a following word for two words
func (l Lexicon) SuggestTwoWordFollower(s string) (string, error) {
	val, ok := l.twoWordFollowers[s]
	if !ok || len(val) == 0 {
		return "", fmt.Errorf("two word follower not found for %s", s)
	}
	fmt.Println(val)
	return val[0], nil
}

func main() {
	r, err := os.Open("data/hhgttg.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	l := NewLexicon()
	l.IngestReader(r)
	w, _ := l.SuggestTwoWordFollower("at a")
	fmt.Println(w)
}
