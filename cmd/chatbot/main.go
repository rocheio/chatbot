package main

import (
	"flag"
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
	bytesRead, err := l.ReadFrom(r)
	if bytesRead == 0 {
		fmt.Println("zero bytes read")
	}
	if err != nil {
		fmt.Println(err)
	}

	// Parse input from command line
	var seed string
	flag.StringVar(&seed, "input", "", "Input text for the chatbot")
	flag.Parse()

	if seed != "" {
		c := l.CommonClause()
		fmt.Println(c.Formatted())
	}
}
