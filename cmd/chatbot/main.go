package main

import (
	"flag"
	"fmt"

	"github.com/rocheio/chatbot/pkg/chat"
)

func main() {
	var input string
	flag.StringVar(&input, "input", "", "Input text for the chatbot")
	flag.Parse()

	if input != "" {
		err := chat.FromInput(input)
		if err != nil {
			fmt.Println("error:", err)
		}
		return
	}

	flag.PrintDefaults()
}
