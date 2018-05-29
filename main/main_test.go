package main

import (
	"strings"
	"testing"
)

func TestTwoWordFollower(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	l := NewLexicon()
	l.IngestString(text)
	expected := "brown"
	actual, err := l.SuggestTwoWordFollower("the quick")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestIngestReader(t *testing.T) {
	r := strings.NewReader(`
		Far out in the uncharted backwaters of the unfashionable
		end  of the  western  spiral arm  of  the Galaxy lies a
		small unregarded yellow sun.
	`)
	l := NewLexicon()
	l.IngestReader(r)
	expected := "of"
	actual, err := l.SuggestTwoWordFollower("unfashionable end")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
