package mux

type set map[string]bool

func newSet(items ...string) set {
	s := set(make(map[string]bool))
	for _, i := range items {
		s.Add(i)
	}
	return s
}

func (s set) Add(item string) {
	s[item] = true
}

func (s set) Has(item string) bool {
	_, ok := s[item]
	return ok
}
