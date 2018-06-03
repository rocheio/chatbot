package chat

import (
	"reflect"
	"strings"
	"testing"
)

func TestOneWordFollower(t *testing.T) {
	l := NewLexicon()
	l.IngestString("a b c")
	expected := "b"
	actual, err := l.Follower("a")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestTwoWordFollower(t *testing.T) {
	type testcase struct {
		corpus    string
		input     string
		expected  string
		shouldErr bool
	}
	testcases := []testcase{
		testcase{"a b c", "a b", "c", false},
		testcase{"a b 1", "a b", "1", false},
		testcase{"w/ Special Chars!", "w special", "chars", false},
		// should use Max count not First
		testcase{"a b c a b z a b z", "a b", "z", false},
		// errors if follower doesn't exist
		testcase{"", "", "", true},
		testcase{"a b", "a b", "", true},
		testcase{"a b", "c d", "", true},
	}
	for _, tc := range testcases {
		l := NewLexicon()
		l.IngestString(tc.corpus)
		actual, err := l.Follower(tc.input)
		if err == nil && tc.shouldErr {
			t.Error("expected error, got nil")
		} else if err != nil && !tc.shouldErr {
			t.Error(err)
		}
		if actual != tc.expected {
			t.Errorf("expected '%s', got '%s'", tc.expected, actual)
		}
	}
}

func TestIngestReader(t *testing.T) {
	r := strings.NewReader(`
		a b.
		a b.
		a c.
	`)
	l := NewLexicon()
	l.IngestReader(r)
	expected := map[string]int{
		"b": 2,
		"c": 1,
	}
	actual := l.oneWordFollowers["a"].m
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestRandomSentence(t *testing.T) {
	l := NewLexicon()
	l.IngestString("a b c d e")
	s := l.RandomSentence("a")
	actual := s.Formatted()
	expected := "A b c d e."
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestCommonClause(t *testing.T) {
	type testcase struct {
		lex      string
		expected string
	}
	testcases := []testcase{
		testcase{"I am robot", "I am robot."},
		testcase{"what is uniquestring", "What is uniquestring."},
	}
	for _, tc := range testcases {
		l := NewLexicon()
		l.IngestString(tc.lex)
		c := l.CommonClause()
		actual := c.Formatted()
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}
