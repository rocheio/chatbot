package chat

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rocheio/chatbot/pkg/lexicon"
)

// Chat is a short session of conversation with a bot
type Chat struct {
	lex      lexicon.Lexicon
	Messages []Message
}

// New returns a new Chat session
func New() Chat {
	l := lexicon.New()
	return Chat{lex: l}
}

// Add appends text from a speaker to it's transcript of Messages
func (c *Chat) Add(speaker, text string) {
	c.Messages = append(c.Messages, Message{text: text, speaker: speaker})
}

// Greet prints a greeting to stdout
func (c *Chat) Greet() {
	speaker := "chatbot"
	greeting := c.lex.RandomGreeting()
	fmt.Printf("%s: %s\n", speaker, greeting)
	c.Add(speaker, greeting)
}

// Respond prints out a response to an input
func (c *Chat) Respond(input string) {
	if c.lex.IsGreeting(input) {
		c.Greet()
		return
	}

	speaker := "chatbot"
	response := "I'm confused"
	fmt.Printf("%s: %s\n", speaker, response)
	c.Add(speaker, response)
}

// Message is a one-direction piece of communication
type Message struct {
	text    string
	speaker string
}

func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.speaker, m.text)
}

// Interactive creates an interactive chat session with a bot
func Interactive() {
	c := New()
	fmt.Println("Starting chat session... (`exit` to exit)")
	c.Greet()
	fmt.Print("you: ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		c.Add("you", text)
		c.Respond(text)
		fmt.Print("you: ")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("saved text: %s\n", c.Messages)
}
