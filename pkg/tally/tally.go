package tally

import (
	"fmt"
	"math/rand"
)

// Returns true if string is in a slice, false otherwise
func contains(l []string, s string) bool {
	for _, item := range l {
		if s == item {
			return true
		}
	}
	return false
}

// Tally tracks the count of strings in a corpus
type Tally struct {
	m map[string]int
}

// New returns a Tally ready to count strings
func New() *Tally {
	return &Tally{make(map[string]int)}
}

// Size returns the number of entries in this Tally
func (t *Tally) Size() int {
	return len(t.m)
}

// Max returns the Tally entry with the highest count
func (t *Tally) Max() string {
	return t.MaxExclude()
}

// MaxExclude returns the Tally with the highest count minus given keys
func (t *Tally) MaxExclude(exclusions ...string) string {
	key := ""
	max := 0
	for k, v := range t.m {
		if contains(exclusions, k) {
			continue
		}
		if v > max {
			key = k
			max = v
		}
	}
	return key
}

// Incr increases the value of a Tally entry
func (t *Tally) Incr(key string) {
	if key == "" {
		return
	}
	t.m[key]++
}

// RandKey returns a random key from the Tally
func (t *Tally) RandKey() string {
	fmt.Println(t.m)
	if len(t.m) == 0 {
		return ""
	}
	r := rand.Intn(len(t.m))
	i := 0
	for k := range t.m {
		if i == r {
			return k
		}
		i++
	}
	return ""
}

// PrintOver writes entries in Tally over n to stdout
func (t *Tally) PrintOver(n int) {
	for k, v := range t.m {
		if v >= n {
			fmt.Println(k, v)
		}
	}
}
