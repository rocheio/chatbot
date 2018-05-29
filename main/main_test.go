package main

import (
	"strings"
	"testing"
)

func TestIngestReader(t *testing.T) {
	r := strings.NewReader("a b c.")
	l := NewLexicon()
	l.IngestReader(r)
	expected := "c"
	actual, err := l.SuggestTwoWordFollower("a b")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

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

// SuggestTwoWordFollower should return based on most popular
func TestTwoWorldFollowerMax(t *testing.T) {
	r := strings.NewReader(`
		Foo bar.
		Foo bar buzz.
		Foo bar baz.
		Foo bar baz.
	`)
	l := NewLexicon()
	l.IngestReader(r)
	expected := "baz"
	actual, err := l.SuggestTwoWordFollower("foo bar")
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
