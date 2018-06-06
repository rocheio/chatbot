package tally

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

// Incr increases the value of a Tally entry
func (t *Tally) Incr(key string) {
	if key == "" {
		return
	}
	t.m[key]++
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
