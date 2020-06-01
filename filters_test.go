package mux

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMethodsFilter(t *testing.T) {
	fil := NewMethodsFilter(http.MethodConnect, http.MethodGet)

	req, err := http.NewRequest(http.MethodGet, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	//-------------------- Another Test Case --------------------
	req, err = http.NewRequest(http.MethodConnect, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	//-------------------- Another Test Case --------------------
	req, err = http.NewRequest(http.MethodDelete, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
}

func TestPathFilter(t *testing.T) {
	fil := NewPathFilter("/{i:int}")

	req, err := http.NewRequest(http.MethodGet, "/32", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
	//-------------------- Another Test Case --------------------
	fil = NewPathFilter("/{s:str}")
	req, err = http.NewRequest(http.MethodGet, "/Viktor", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/$32", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
	//-------------------- Another Test Case --------------------
	fil = NewPathFilter("/p/{name:str}/{age:int}")
	req, err = http.NewRequest(http.MethodGet, "/p/Alex/42", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/p/32/Alex", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
	//-------------------- Another Test Case --------------------
	fil = NewPathFilter("/p/{age:nat}")
	req, err = http.NewRequest(http.MethodGet, "/p/42", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/p/-32", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
	//-------------------- Another Test Case --------------------
	fil = NewPathFilter("/pub/.*")
	req, err = http.NewRequest(http.MethodGet, "/pub/lisn/index.html", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/p/-32", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
	//-------------------- Another Test Case --------------------
	fil = NewPathFilter(`/pub/fail/{file:\d{3}\.html}`)
	req, err = http.NewRequest(http.MethodGet, "/pub/fail/404.html", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "/pub/fail/a404.html", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the PathFilter matched an incorrect path")
	}
}

func TestPathFilterVars(t *testing.T) {
	rtr := New().Path("/r/{article:str}/{id:nat}").HandleFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars, ok := Vars(r)
			if !ok {
				t.Error("the Vars function failed to retreive path variables")
			}
			article := vars["article"]
			id := vars["id"]
			s := fmt.Sprintf("#%d - %s", id, article)
			if s != "#42 - Computers" {
				t.Errorf("got '%s'; expected '#42 - Computers'", s)
			}
		},
	)

	rec, req, err := request(http.MethodGet, "/r/Computers/42", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	rtr.ServeHTTP(rec, req)
	//-------------------- Another Test Case --------------------
	rtr.Path(`/r/{article:str}/{id:\w\d}`).HandleFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars, ok := Vars(r)
			if !ok {
				t.Error("the Vars function failed to retreive path variables")
			}
			article := vars["article"]
			id := vars["id"]
			s := fmt.Sprintf("%s @ %s", id, article)
			if s != "a2 @ Computers" {
				t.Errorf("got '%s'; expected 'a2 @ Computers'", s)
			}
		},
	)

	rec, req, err = request(http.MethodGet, "/r/Computers/a2", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	rtr.ServeHTTP(rec, req)
}

func TestPathPrefixFilter(t *testing.T) {
	api := New().PathPrefix("/api")
	api.Subrouter().Path("/song/{id:int}").HandleFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vars, ok := Vars(r)
			if !ok {
				t.Errorf("the Vars function failed to retrieve variables")
			}
			id := vars["id"]
			s := fmt.Sprintf("Song #%d", id)
			if s != "Song #42" {
				t.Errorf("got '%s'; expected 'Song #42'", s)
			}
		},
	)

	rec, req, err := request(http.MethodGet, "/api/song/42", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	api.ServeHTTP(rec, req)
}

func TestSchemes(t *testing.T) {
	fil := NewSchemesFilter("http")

	req, err := http.NewRequest(http.MethodGet, "http://foo.com/api", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the SchemesFilter did not match a correct path")
	}
	req, err = http.NewRequest(http.MethodGet, "https://foo.com/api", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if fil.Match(req) {
		t.Error("the SchemesFilter matched an incorrect path")
	}
}
