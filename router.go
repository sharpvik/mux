package mux

import (
	"context"
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

	// Preview is a function of type View that is called whenever current
	// Router's ServeHTTP method is triggered unless that Router doesn't have a
	// Preview. It is initially set to nil on Router creation, however, you can
	// put here any function you wish to *always* be executed on ServeHTTP.
	//
	// NOTICE: Preview is executed even if current Router is not the final
	// handler for this particular request.
	Preview View

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
		nil, // Preview
		nil, // View
		DefaultFailMessage,
		nil, // routes
		NewFilters(),
	}
}

// ServeHTTP method is here in order to ensure that Router implements the
// http.Handler interface. It is invoked automatically by http.Server if you
// assign Router in question as server's Handler.
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Cut path prefix (if set) from the reuqest URL path.
	if rtr.filters.PathPrefix != nil {
		r.URL.Path = strings.TrimPrefix(
			r.URL.Path, string(*rtr.filters.PathPrefix),
		)
	}

	// Parse path variables and alter http.Request.Context.
	r = rtr.vars(r)

	// Must call Preview if present.
	if rtr.Preview != nil {
		rtr.Preview(w, r, rtr.Context)
	}

	// 1. Check if there are routes with matching filters.
	// 2. If not, call View if present.
	// 3. If everything else failed, respond with a Fail message.
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
	// Create new Router that inherits its parent's Context.
	sub := New(rtr.Context)

	// Add it to parent's routes.
	rtr.routes = append(rtr.routes, sub)

	return sub
}

// Methods returns pointer to the same rtr instance while altering its methods
// filter.
//
// NOTICE: If methods filter has already been set for this Router instance, it
// will get replaced!
func (rtr *Router) Methods(methods ...string) *Router {
	rtr.filters.Methods = NewMethodsFilter(methods...)
	return rtr
}

// Path returns pointer to the same rtr instance while altering its path filter.
//
// NOTICE: This method replaces router's PathFilter with a newly created
// instance while setting PathPrefix to nil.
func (rtr *Router) Path(path string) *Router {
	rtr.filters.Path = NewPathFilter(path)
	rtr.filters.PathPrefix = nil
	return rtr
}

// PathPrefix returns pointer to the same rtr instance while altering its path
// prefix filter.
//
// NOTICE: This method replaces router's PathPrefixFilter with a newly created
// instance while setting PathFilter to nil.
func (rtr *Router) PathPrefix(prefix string) *Router {
	rtr.filters.PathPrefix = NewPathPrefixFilter(prefix)
	rtr.filters.Path = nil
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

// vars method parses variables from request using the PathFilter.Path and
// stores them in http.Request.Context.
//
// This is a non-exported method that's only triggered by Router's ServeHTTP
// method. Therefore, we can assume that the Request given to us matches all
// Router's filters including the PathFilter (if present).
func (rtr *Router) vars(r *http.Request) *http.Request {
	pathfil := rtr.filters.Path

	// Check if PathFilter is present.
	if pathfil == nil {
		return r
	}

	// Check if PathFilter has variables.
	if !pathfil.hasVars {
		return r
	}

	// At this point, we know that rtr has a PathFilter with vars.
	vars := make(map[string]interface{})
	path := pathfil.Path

	// Slicing the first element away because it is always going to be an empty
	// string since the first character is always a slash.
	fsplit := strings.Split(path, "/")[1:]
	rsplit := strings.Split(r.URL.Path, "/")[1:]

	// Linear pattern matching. The pat here is a field from the filter path,
	// exp is a request path field we want to match towards. Both are strings.
	// For example, pat = "{n:int}"; exp = "42".
	for i, pat := range fsplit {
		exp := rsplit[i]

		// Skip all patterns that are not variables. No need to validate them.
		if !isVar(pat) {
			continue
		}

		name, typ := varData(pat)

		switch typ {
		case "int":
			// Discarding the error here because we know for sure that exp
			// passed regex test for number.
			vars[name], _ = strconv.Atoi(exp)

		case "nat":
			vars[name], _ = strconv.ParseUint(exp, 10, 0)

		case "str":
			vars[name] = exp

		default: // regex type
			vars[name] = exp
		}
	}

	return r.WithContext(context.WithValue(r.Context(), varsKey, vars))
}
