package mux

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	routes []*Router

	// filters is a set of filters that are used to check whether this Router
	// instance should be used for the request at hand.
	filters *Filters
}

// DefaultFailMessage is just a string with some dummy failure message.
const DefaultFailMessage = "Handler node did not have a view assigned to it."

// New is a constructor used to create the root of a routing tree. Root doesn't
// need any filters as it is invoked automatically by the server anyway.
// The routes will be added later, using Router's methods.
func New(ctx Context) *Router {
	return &Router{
		ctx,
		nil,
		DefaultFailMessage,
		nil,
		NewFilters(),
	}
}

// ServeHTTP method is here in order to ensure that Router implements the
// http.Handler interface. It is invoked automatically by http.Server if you set
// its Handler to the Router in question.
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sub, match := rtr.Match(r); match {
		sub.ServeHTTP(w, r)
	} else if rtr.View != nil {
		rtr.View(w, r, rtr.Context)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, rtr.Fail)
	}
}

// Subrouter method returns pointer to a new sub-router instance that inherits
// context from its parents.
func (rtr *Router) Subrouter() *Router {
	// Create new Router with the same Context.
	sub := New(rtr.Context)

	// Add it to parent's routes.
	rtr.routes = append(rtr.routes, sub)

	return sub
}

// Methods returns pointer to the same rtr instance while altering its filters.
func (rtr *Router) Methods(methods ...string) *Router {
	rtr.filters.Methods = NewMethodsFilter(methods...)
	return rtr
}

// Path returns pointer to the same rtr instance while altering its filters.
func (rtr *Router) Path(path string) *Router {
	rtr.filters.Path = NewPathFilter(path)
	return rtr
}

// Match method must go through all registered routes one by one and check if
// their filters match the request. It returns the first sub-router where
// filters matched and a boolean value indicating that there was a match.
// If there was no match, it returns nil as the sub-router while setting the
// second value to false.
func (rtr *Router) Match(r *http.Request) (sub *Router, match bool) {
	for _, route := range rtr.routes {
		if route.filters.Match(r) {
			return route, true
		}
	}
	return nil, false
}

// Vars method parses variables from request using the PathFilter.Path.
func (rtr *Router) Vars(r *http.Request) map[string]interface{} {
	empty := make(map[string]interface{})
	vars := make(map[string]interface{})

	path := rtr.filters.Path.Path

	// Slicing the first element away because it is always going to be an empty
	// string since the first character is always a slash.
	fsplit := strings.Split(path, "/")[1:]
	rsplit := strings.Split(r.URL.Path, "/")[1:]

	if len(fsplit) != len(rsplit) {
		return empty
	}

	// Linear pattern matching. The pat here is a field from the filter path,
	// exp is a request path field we want to match towards. Both are strings.
	// For example, pat = "{n:int}"; exp = "42".
	for i, pat := range fsplit {
		exp := rsplit[i]

		if isVar(pat) {
			name, typ := varData(pat)

			switch typ {
			case Int:
				val, err := strconv.Atoi(exp)
				if err != nil {
					return empty
				}
				vars[name] = val

			case Str:
				val := exp
				vars[name] = val
			}
		} else if pat != exp {
			return empty
		}
	}

	return vars
}
