package chat

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"unicode"
)

var knownSubjects = []string{"i", "you", "he", "she", "it"}
var knownVerbs = []string{"am", "is", "like", "want"}
var knownPredicates = []string{"me", "you", "he", "she", "it", "robot"}

// normalize returns a string in lowercase with only spaces for whitespace
func normalize(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return alphanum(strings.ToLower(s))
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

// Return string with first char made uppercase
func sentenceCase(s string) string {
	if len(s) == 0 {
		return ""
	} else if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// Returns true if string is in a slice, false otherwise
func contains(l []string, s string) bool {
	for _, item := range l {
		if s == item {
			return true
		}
	}
	return false
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
func (t *Tally) Size() int {
	return len(t.m)
}

// Max returns the Tally entry with the highest count
func (t *Tally) Max() string {
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
func (t *Tally) Incr(key string) {
	t.m[key]++
}

// RandKey returns a random key from the Tally
func (t *Tally) RandKey() string {
	fmt.Println(t.m)
	if len(t.m) == 0 {
		return ""
	}
	r := rand.Intn(len(t.m))
	i := 0
	for k := range t.m {
		if i == r {
			return k
		}
		i++
	}
	return ""
}

// Lexicon defines a vocabulary for text structures
type Lexicon struct {
	// maps known chains of words to number of occurances
	oneWordFollowers map[string]*Tally
	twoWordFollowers map[string]*Tally
	// maps for existence of types of words to number of occurances
	subjects   *Tally
	verbs      *Tally
	predicates *Tally
}

// NewLexicon returns an empty Lexicon ready to ingest data
func NewLexicon() Lexicon {
	return Lexicon{
		oneWordFollowers: make(map[string]*Tally),
		twoWordFollowers: make(map[string]*Tally),
		subjects:         NewTally(),
		verbs:            NewTally(),
		predicates:       NewTally(),
	}
}

// IngestString adds a strings contents to this Lexicon
func (l Lexicon) IngestString(s string) {
	s = normalize(s)
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
	// Assume all three word sentences are simple clauses
	// TODO -- refine this to work with 4-words using `a` / `the`
	// TODO -- test with large corpuses and add filters for false positives
	if len(parts) == 3 {
		l.subjects.Incr(parts[0])
		l.verbs.Incr(parts[1])
		l.predicates.Incr(parts[2])
		return
	}
	// hacky way to frontload common words for testing other components
	// TODO -- get better logic to learn from and ditch this
	for _, p := range parts {
		if contains(knownSubjects, p) {
			l.subjects.Incr(p)
		}
		if contains(knownVerbs, p) {
			l.verbs.Incr(p)
		}
		if contains(knownPredicates, p) {
			l.predicates.Incr(p)
		}
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

// Follower returns a following word for a string
func (l Lexicon) Follower(s string) (string, error) {
	var tally *Tally
	var ok bool
	strLen := len(strings.Split(s, " "))
	if strLen == 1 {
		tally, ok = l.oneWordFollowers[s]
	} else if strLen == 2 {
		tally, ok = l.twoWordFollowers[s]
	}
	if !ok || tally.Size() == 0 {
		return "", fmt.Errorf("follower not found for '%s'", s)
	}
	return tally.Max(), nil
}

// TryAllFollowers looks for all permutations of a Sentence for followers.
// Followers with a match on more words take precedence.
func (l Lexicon) TryAllFollowers(s Sentence) (string, error) {
	maxDist := 2
	if len(s.words) < maxDist {
		maxDist = len(s.words)
	}
	for i := maxDist; i > 0; i-- {
		word, err := l.Follower(strings.Join(s.Last(i), " "))
		if err == nil {
			return word, nil
		}
	}
	return "", fmt.Errorf("no matches found for sentence %v", s.words)
}

// RandomSentence returns a sentence from a seed word
func (l Lexicon) RandomSentence(word string) Sentence {
	s := NewSentence()
	for {
		err := s.Add(word)
		if err != nil {
			break
		}

		word, err = l.TryAllFollowers(s)
		if err != nil {
			break
		}
	}
	return s
}

// CommonClause returns the most common simple clause in the Lexicon
func (l Lexicon) CommonClause() Clause {
	return Clause{
		Subject:   l.subjects.Max(),
		Verb:      l.verbs.Max(),
		Predicate: l.predicates.Max(),
	}
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
	if n > len(s.words) {
		return s.words
	}
	return s.words[len(s.words)-n:]
}

// Formatted returns the sentence as a properly formatted string
func (s *Sentence) Formatted() string {
	return sentenceCase(strings.Join(s.words, " ")) + "."
}

// Clause is the smallest grammatical structure for a proposition
type Clause struct {
	Subject   string
	Verb      string
	Predicate string
}

// Formatted returns the clause as a single formatted string
func (c *Clause) Formatted() string {
	return sentenceCase(
		fmt.Sprintf("%s %s %s", c.Subject, c.Verb, c.Predicate),
	) + "."
}
