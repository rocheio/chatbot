package main

import (
	"fmt"
	"os"

	"github.com/rocheio/chatbot/pkg/chat"
)

func main() {
	l := chat.NewLexicon()

	r, err := os.Open("data/hhgttg.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l.IngestReader(r)

	s := l.RandomSentence("hello")
	fmt.Println(s.Formatted())

	c := l.CommonClause()
	fmt.Println(c.Formatted())
}
