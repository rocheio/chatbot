package main

import (
	"fmt"
	"os"

	"github.com/rocheio/chatbot/pkg/chat"
)

func main() {
	// Build a lexicon from Hitchiker's Guide to the Galaxy
	l := chat.NewLexicon()
	r, err := os.Open("data/hhgttg.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l.ReadReader(r)

	// Print out a simple Markov chain sentence
	s := l.RandomSentence("hello")
	fmt.Println("\nSentence:\n" + s.Formatted() + "\n")

	// Print out a simple clause
	c := l.CommonClause()
	fmt.Println("Simple Clause:\n" + c.Formatted() + "\n")

	// Print out a progression of unique simple clauses
	fmt.Println("Progressive Clauses:")
	var excluded []string
	for i := 0; i < 5; i++ {
		c := l.ExclusionClause(excluded)
		fmt.Println(c.Formatted())
		excluded = append(excluded, c.Subject)
		excluded = append(excluded, c.Predicate)
		excluded = append(excluded, c.Verb)
	}
}
