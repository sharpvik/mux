package mux

import (
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
