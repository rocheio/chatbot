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
