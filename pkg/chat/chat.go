package chat

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/rocheio/chatbot/pkg/tally"
)

var definiteArticle = "the"
var negativeArticle = "no"
var indefiniteArticles = []string{"a", "an", "some"}
var nouns = []string{"me", "you", "he", "she", "it"}
var verbs = []string{"am", "is", "was", "has", "wants"}

// normalizeWhitespace returns a string separated by only single spaces
func normalizeWhitespace(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
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

// Removes a string from a slice
func removeStr(l []string, s string) []string {
	for i, val := range l {
		if val == s {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// Lexicon defines a vocabulary for text structures
type Lexicon struct {
	// maps known chains of words to number of occurances
	oneWordFollowers map[string]*tally.Tally
	twoWordFollowers map[string]*tally.Tally
	// maps for existence of types of words to number of occurances
	nouns *tally.Tally
	verbs *tally.Tally
}

// NewLexicon returns an empty Lexicon ready to ingest data
func NewLexicon() Lexicon {
	return Lexicon{
		oneWordFollowers: make(map[string]*tally.Tally),
		twoWordFollowers: make(map[string]*tally.Tally),
		nouns:            tally.New(),
		verbs:            tally.New(),
	}
}

// ReadString adds a strings contents to this Lexicon
func (l Lexicon) ReadString(s string) {
	s = alphanum(strings.ToLower(s))
	parts := strings.Split(s, " ")
	for i := 0; i < len(parts)-1; i++ {
		key := parts[i]
		val := parts[i+1]
		if l.oneWordFollowers[key] == nil {
			l.oneWordFollowers[key] = tally.New()
		}
		l.oneWordFollowers[key].Incr(val)
	}
	for i := 0; i < len(parts)-2; i++ {
		key := strings.Join([]string{parts[i], parts[i+1]}, " ")
		val := parts[i+2]
		if l.twoWordFollowers[key] == nil {
			l.twoWordFollowers[key] = tally.New()
		}
		l.twoWordFollowers[key].Incr(val)
	}
	// Strip articles to find potential simple clauses
	if len(parts) >= 3 && len(parts) <= 5 {
		for _, p := range parts {
			if p == definiteArticle ||
				p == negativeArticle ||
				contains(indefiniteArticles, p) {
				parts = removeStr(parts, p)
			}
		}
	}
	// Assume all 3 word sentences at this point are simple clauses
	// TODO -- test with large corpuses and add filters for false positives
	if len(parts) == 3 {
		l.nouns.Incr(parts[0])
		l.verbs.Incr(parts[1])
		l.nouns.Incr(parts[2])
		return
	}
	// hacky way to frontload common words for testing other components
	// TODO -- get better logic to learn from and ditch this
	for _, p := range parts {
		if contains(nouns, p) {
			l.nouns.Incr(p)
		}
		if contains(verbs, p) {
			l.verbs.Incr(p)
		}
	}
}

// ReadSentence adds a sentence to this Lexicon
func (l Lexicon) ReadSentence(s string) {
	s = normalizeWhitespace(s)
	fmt.Println("\n", s)
	l.ReadString(s)
}

// ReadReader adds all text from a Reader to this Lexicon
func (l Lexicon) ReadReader(r io.Reader) {
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
		l.ReadSentence(sentence)
	}
}

// Follower returns a following word for a string
func (l Lexicon) Follower(s string) (string, error) {
	var t *tally.Tally
	var ok bool
	strLen := len(strings.Split(s, " "))
	if strLen == 1 {
		t, ok = l.oneWordFollowers[s]
	} else if strLen == 2 {
		t, ok = l.twoWordFollowers[s]
	}
	if !ok || t.Size() == 0 {
		return "", fmt.Errorf("follower not found for '%s'", s)
	}
	return t.Max(), nil
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

// CommonClause returns the most common simple clause in l
func (l Lexicon) CommonClause() Clause {
	return l.ExclusionClause(nil)
}

// ExclusionClause returns a simple clause from l with items excluded
func (l Lexicon) ExclusionClause(excluded []string) Clause {
	subj := l.nouns.MaxExclude(excluded...)
	excluded = append(excluded, subj)
	verb := l.verbs.MaxExclude(excluded...)
	pred := l.nouns.MaxExclude(excluded...)
	return Clause{
		Subject:   subj,
		Verb:      verb,
		Predicate: pred,
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
