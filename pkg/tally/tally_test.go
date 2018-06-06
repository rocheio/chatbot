package tally

import (
	"testing"
)

func TestSize(t *testing.T) {
	tally := New()
	actual := tally.Size()
	expected := 0
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestMax(t *testing.T) {
	tally := New()
	actual := tally.Max()
	expected := ""
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestMaxExclude(t *testing.T) {
	tally := New()
	tally.Incr("a")
	tally.Incr("a")
	tally.Incr("b")
	actual := tally.MaxExclude("a")
	expected := "b"
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestIncr(t *testing.T) {
	tally := New()
	tally.Incr("foo")

	actual := tally.m["foo"]
	expected := 1
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	// don't include blanks in Tallys because they keep
	// slipping in from parsers and it's a pain to keep checking
	// this logic in higher-level code
	tally.Incr("")
	actual = tally.m[""]
	expected = 0
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
