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
		chat.Interactive()
		return
	}

	if input != "" {
		c := chat.New()
		response, err := c.Response(input)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(response)
		return
	}

	flag.PrintDefaults()
}
