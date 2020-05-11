package mux

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMethodsFilter(t *testing.T) {
	root := New(Cont{"laughing out loud"})

	sub := root.Subrouter().Methods(http.MethodGet, http.MethodDelete)
	sub.View = func(w http.ResponseWriter, r *http.Request, ctx Context) {
		fmt.Fprintf(w, "Method: '%s'", r.Method)
	}

	rec, req, err := request(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}

	err = result(
		root, rec, req,
		func(res *http.Response) error {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf(
					"status in response: '%v'; expected '200 OK'",
					res.Status,
				)
			}

			if body, _ := ioutil.ReadAll(res.Body); string(body) != "Method: 'GET'" {
				return fmt.Errorf(
					"response body: %s; expected: `Method: 'GET'`",
					body,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestPathFilter(t *testing.T) {
	root := New(Cont{"laughing out loud"})

	sub := root.Subrouter().Path("/lol")
	sub.View = func(w http.ResponseWriter, r *http.Request, ctx Context) {
		fmt.Fprintf(w, "lol")
	}

	rec, req, err := request(http.MethodGet, "/lol", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}

	err = result(
		root, rec, req,
		func(res *http.Response) error {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf(
					"status in response: '%v'; expected '200 OK'",
					res.Status,
				)
			}

			if body, _ := ioutil.ReadAll(res.Body); string(body) != "lol" {
				return fmt.Errorf(
					"response body: %s; expected: 'lol'",
					body,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Error(err)
	}
}
