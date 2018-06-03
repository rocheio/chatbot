package tally

import (
	"testing"
)

func TestRandKey(t *testing.T) {
	tally := New()
	tally.Incr("a")
	expected := "a"
	actual := tally.RandKey()
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
