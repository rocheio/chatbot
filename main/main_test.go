package main

import (
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
