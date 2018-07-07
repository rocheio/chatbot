package main

import (
	"flag"
	"fmt"

	"github.com/rocheio/chatbot/pkg/chat"
)

func main() {
	var interactive bool
	var input string

	flag.BoolVar(&interactive, "interactive", false, "Interactive chatbot session")
	flag.StringVar(&input, "input", "", "Input text for the chatbot")
	flag.Parse()

	if interactive {
		err := chat.Interactive()
		if err != nil {
			fmt.Println("error:", err)
		}
		return
	}

	if input != "" {
		c := chat.New()
		c.Respond(input)
		return
	}

	flag.PrintDefaults()
}
