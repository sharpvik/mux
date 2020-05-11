package mux

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Cont struct {
	msg string
}

func TestRootRouter(t *testing.T) {
	root := New(Cont{"laughing out loud"})
	root.Fail = "lol fail"

	// View function not set.
	rec, req, err := request(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}

	err = result(
		root, rec, req,
		func(res *http.Response) error {
			if res.StatusCode != http.StatusNotImplemented {
				return fmt.Errorf(
					"status in response: '%v'; expected '501 Not Implemented'",
					res.Status,
				)
			}

			if body, _ := ioutil.ReadAll(res.Body); string(body) != "lol fail" {
				return fmt.Errorf(
					"response body: %s; expected 'lol fail'",
					body,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Error(err)
	}

	// After setting the View.
	root.View = func(w http.ResponseWriter, r *http.Request, ctx Context) {
		fmt.Fprint(w, ctx.(Cont).msg)
	}

	rec, req, err = request(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}

	err = result(
		root, rec, req,
		func(res *http.Response) error {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf(
					"status in response: '%v'; expected '200 OK'", res.Status,
				)
			}

			if body, _ := ioutil.ReadAll(res.Body); string(body) != "laughing out loud" {
				return fmt.Errorf(
					"response body: %s; expected 'laughing out loud'",
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

func request(method string, addr string, body io.Reader) (
	w *httptest.ResponseRecorder, r *http.Request, err error,
) {
	r, err = http.NewRequest(method, addr, body)
	if err != nil {
		return
	}
	w = httptest.NewRecorder()
	return
}

func result(
	rtr *Router,
	rec *httptest.ResponseRecorder,
	req *http.Request,
	fun func(*http.Response) error,
) error {
	rtr.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	return fun(res)
}
