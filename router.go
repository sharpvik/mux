package mux

import (
	"fmt"
	"net/http"
)

// Router represents the node of a routing tree.
type Router struct {
	// Context contains all of the things you'd like to be exposed to the View
	// handler function. Its type is an empty interface so that you could
	// declare your own custom context type and use it here.
	Context Context

	// View is a handler function that is triggered if current request did not
	// match any filters. It may hold an actual handler function, for example,
	// if this Router instance is the leaf node of the routing tree.
	// Alternatively, it may hold a fail handler function of your choice.
	View View

	// Fail is a failure message written to http.ResponseWriter by the ServeHTTP
	// method in case current request did not match any routes and the View
	// handler function was not set.
	//
	// Initially its value is set to be DefaultFailMessage, but you can easily
	// change it if you want.
	Fail string

	// routes is a slice of sub-routers.
	routes []Router

	// filters is a slice of filters that are used to check whether this Router
	// instance should be used for the request at hand.
	filters []Filter
}

// DefaultFailMessage is just a string with some dummy failure message.
const DefaultFailMessage = "Handler node did not have a view assigned to it."

// RootRouter is a constructor used to create the root of a routing tree. Root
// doesn't need any filters as it is invoked automatically by the server anyway.
// The routes will be added later, using Router's methods.
func RootRouter(ctx Context) *Router {
	return &Router{
		ctx,
		nil,
		DefaultFailMessage,
		nil,
		nil,
	}
}

// ServeHTTP method is here in order to ensure that Router implements the
// http.Handler interface. It is invoked automatically by http.Server if you set
// its Handler to the Router in question.
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sub, match := rtr.match(r); match {
		sub.ServeHTTP(w, r)
	} else if rtr.View != nil {
		rtr.View(w, r, rtr.Context)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, rtr.Fail)
	}
}

// match method must go through all registered routes one by one and check if
// their filters match the request. It returns the first sub-router where
// filters matched and a boolean value indicating that there was a match.
// If there was no match, it returns itself as the sub-router while setting
// second value to false.
func (rtr *Router) match(r *http.Request) (sub *Router, match bool) {
	return rtr, false
}
