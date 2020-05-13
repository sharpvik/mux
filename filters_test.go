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

	req, err = http.NewRequest(http.MethodConnect, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	if !fil.Match(req) {
		t.Error("the PathFilter did not match a correct path")
	}

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
}

func TestPathFilterVars(t *testing.T) {
	rtr := New(&Cont{"lol"}).Path("/r/{article:str}/{id:int}")
	rtr.View = func(w http.ResponseWriter, r *http.Request, ctx Context) {
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
	}

	rec, req, err := request(http.MethodGet, "/r/Computers/42", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	rtr.ServeHTTP(rec, req)
}
