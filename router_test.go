package mux

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootRouter(t *testing.T) {
	root := New().FailFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "lol fail")
	})

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
	//-------------------- Another Test Case --------------------
	// After setting the View.
	root.HandleFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "laughing out loud")
		},
	)

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

func TestRouterMiddleware(t *testing.T) {
	rtr := New().
		UseFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("middleware", "ok")
		}).
		HandleFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "middleware ok")
		})
	rec, req, err := request(http.MethodGet, "/", nil)
	assert.NoError(t, err, "request failed:", err)
	err = result(rtr, rec, req,
		func(r *http.Response) (err error) {
			if ok := r.Header.Get("middleware"); ok != "ok" {
				return errors.New("middleware did not work")
			}
			var body []byte
			if body, _ = ioutil.ReadAll(r.Body); string(body) != "middleware ok" {
				return errors.New("handler was ignored after middleware application")
			}
			return
		})
	assert.NoError(t, err, "middleware failed:", err)
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
