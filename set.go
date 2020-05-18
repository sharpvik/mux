package mux

// set is a map-based data structure that allows us to add and check presence
// of any string in constant time.
type set map[string]bool

// newSet returns a new set instance with all items added to it.
func newSet(items ...string) set {
	s := set(make(map[string]bool))
	for _, i := range items {
		s.Add(i)
	}
	return s
}

// Add method accepts an item you want to add to the set.
func (s set) Add(item string) {
	s[item] = true
}

// Has method returns a boolean flag that tells you whether accepted string
// has been previously added to this set.
func (s set) Has(item string) bool {
	_, ok := s[item]
	return ok
}
