package mux

import "net/http"

// View represents the default handler function type.
type View func(http.ResponseWriter, *http.Request)

// ServeHTTP method ensures that View implements http.Handler interface.
func (v View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v(w, r)
}

// contextKey is an alias for int that we use as a custom type for request
// context key.
type contextKey int

// varsKey is a context key for request variables.
const varsKey contextKey = iota
