package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const maxSentenceLength = 10

// normalize returns a string in lowercase with only spaces for whitespace
func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// alphanum returns a string with only alphanumeric characters
func alphanum(s string) string {
	var clean string
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			clean += string(r)
		}
	}
	return clean
}

// Tally tracks the count of strings in a corpus
type Tally struct {
	m map[string]int
}

// NewTally returns a Tally ready to count strings
func NewTally() *Tally {
	return &Tally{make(map[string]int)}
}

// Size returns the number of entries in this Tally
func (t Tally) Size() int {
	return len(t.m)
}

// Max returns the Tally entry with the highest count
func (t Tally) Max() string {
	key := ""
	max := 0
	for k, v := range t.m {
		if v > max {
			key = k
			max = v
		}
	}
	return key
}

// Incr increases the value of a Tally entry
func (t Tally) Incr(key string) {
	t.m[key]++
}

// Lexicon defines a vocabulary for text structures
type Lexicon struct {
	oneWordFollowers map[string]*Tally
	twoWordFollowers map[string]*Tally
}

// NewLexicon returns a Lexicon ready to ingest data
func NewLexicon() Lexicon {
	return Lexicon{
		oneWordFollowers: make(map[string]*Tally),
		twoWordFollowers: make(map[string]*Tally),
	}
}

// IngestString adds a strings contents to this Lexicon
func (l Lexicon) IngestString(s string) {
	s = alphanum(normalize(s))
	parts := strings.Split(s, " ")
	for i := 0; i < len(parts)-1; i++ {
		key := parts[i]
		val := parts[i+1]
		if l.oneWordFollowers[key] == nil {
			l.oneWordFollowers[key] = NewTally()
		}
		l.oneWordFollowers[key].Incr(val)
	}
	for i := 0; i < len(parts)-2; i++ {
		key := strings.Join([]string{parts[i], parts[i+1]}, " ")
		val := parts[i+2]
		if l.twoWordFollowers[key] == nil {
			l.twoWordFollowers[key] = NewTally()
		}
		l.twoWordFollowers[key].Incr(val)
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
		l.IngestString(sentence)
	}
}

// OneWordFollower returns a following word for a single word
func (l Lexicon) OneWordFollower(s string) (string, error) {
	tally, ok := l.oneWordFollowers[s]
	if !ok || tally.Size() == 0 {
		return "", fmt.Errorf("one word follower not found for %s", s)
	}
	return tally.Max(), nil
}

// TwoWordFollower returns a following word for two words
func (l Lexicon) TwoWordFollower(s string) (string, error) {
	tally, ok := l.twoWordFollowers[s]
	if !ok || tally.Size() == 0 {
		return "", fmt.Errorf("two word follower not found for '%s'", s)
	}
	return tally.Max(), nil
}

// RandomSentence returns a sentence from a seed string
func (l Lexicon) RandomSentence(seed string) string {
	var sentence []string
	var err error
	word := seed

	for {
		if len(sentence) == maxSentenceLength {
			break
		}
		sentence = append(sentence, word)

		if len(sentence) > 2 {
			word, err = l.TwoWordFollower(
				strings.Join(sentence[len(sentence)-2:], " "),
			)
			if err == nil {
				continue
			}
		}

		word, err = l.OneWordFollower(sentence[len(sentence)-1])
		if err == nil {
			continue
		}

		break
	}
	return strings.Join(sentence, " ")
}

func main() {
	l := NewLexicon()

	r, err := os.Open("data/hhgttg.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l.IngestReader(r)

	fmt.Println(l.RandomSentence("hello"))
}
