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
	greeting := "hello"
	c.Add(speaker, greeting)
	fmt.Printf("%s: %s\n", speaker, greeting)
}

// Response prints out a response to an input
func (c *Chat) Response(input string) (string, error) {
	if c.lex.IsGreeting(input) {
		return "", nil
	}

	return "I'm confused", nil
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
		c.Greet()
		fmt.Print("you: ")
	}

	if scanner.Err() != nil {
		// handle error.
	}

	fmt.Printf("saved text: %s\n", c.Messages)
}
