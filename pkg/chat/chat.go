package chat

import (
	"fmt"
	"os"

	"github.com/rocheio/chatbot/pkg/lexicon"
)

// FromInput prints out a response to an input string
func FromInput(input string) error {
	// Build a lexicon from Hitchiker's Guide to the Galaxy
	l := lexicon.New()
	r, err := os.Open("data/hhgttg.txt")
	if err != nil {
		return err
	}
	bytesRead, err := l.ReadFrom(r)
	if bytesRead == 0 {
		return fmt.Errorf("zero bytes read from '%s'", r.Name())
	}
	if err != nil {
		fmt.Println(err)
	}

	c := l.CommonClause()
	fmt.Println(c.Formatted())

	return nil
}
