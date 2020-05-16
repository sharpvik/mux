package mux

import "testing"

func TestSet(t *testing.T) {
	s := newSet()
	s.Add("GET")
	if !s.Has("GET") {
		t.Errorf("set failed to add item")
	}
	if s.Has("POST") {
		t.Errorf("set claims to have item that hasn't been added")
	}
}
