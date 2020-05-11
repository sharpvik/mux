package mux

import "net/http"

// Context type is an empty interface that represents all the things you wish to
// expose to your Views (handler functions). Feel free to declare your own
// context types and use them with the router.
type Context interface{}

// View is a special function type that represents a handler function. The last
// parameter it expects, represents context you wish to expose to the View. It
// is an empty interface in order to allow you to create your own context type.
type View func(http.ResponseWriter, *http.Request, Context)

// Filter is an interface type that represents essential functionality of a
// filter.
type Filter interface {
	Match(r *http.Request) bool
}
