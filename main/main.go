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

// RandomSentence returns a sentence from a seed word
func (l Lexicon) RandomSentence(word string) string {
	s := NewSentence()

	for {
		err := s.Add(word)
		if err != nil {
			break
		}

		if len(s.words) > 2 {
			word, err = l.TwoWordFollower(
				strings.Join(s.Last(2), " "),
			)
			if err == nil {
				continue
			}
		}

		word, err = l.OneWordFollower(s.Last(1)[0])
		if err == nil {
			continue
		}

		break
	}
	return s.Formatted()
}

// Sentence represents a grammatically correct series of words
type Sentence struct {
	words     []string
	maxLength int
}

// NewSentence returns a new empty sentence with default requirements
func NewSentence() Sentence {
	return Sentence{
		words:     []string{},
		maxLength: 10,
	}
}

// Add adds a word to this sentence
func (s *Sentence) Add(w string) error {
	if len(s.words) > s.maxLength {
		return fmt.Errorf("sentence is already at maxlength %d", s.maxLength)
	}
	s.words = append(s.words, w)
	return nil
}

// Last returns the last n words from this Sentence
func (s *Sentence) Last(n int) []string {
	return s.words[len(s.words)-n:]
}

// Formatted returns the sentence as a properly formatted string
func (s *Sentence) Formatted() string {
	str := strings.Join(s.words, " ")
	str = strings.ToUpper(string(str[0])) + str[1:] + "."
	return str
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
