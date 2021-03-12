package mux

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Router represents the node of a routing tree.
type Router struct {
	handler http.Handler

	// Fail is a failure message written to http.ResponseWriter by the ServeHTTP
	// method in case current request did not match any routes and the View
	// handler function was not set.
	//
	// Initially its value is set to be DefaultFailMessage, but you can easily
	// change it if you want.
	fail http.Handler

	// routes is a slice of sub-routers.
	routes []*Router

	// filters is a set of filters that are used to check whether this Router
	// instance should be used for the request at hand.
	filters *Filters

	// middleware is just a list of handlers that are applied to the request
	// before it is passed to the final Router's handler or a subroute.
	middleware []http.Handler
}

// DefaultFailHandler is a default handler attached to every Router. Use
// Router.Fail to specify a custom one.
var DefaultFailHandler = http.NotFoundHandler()

// New is a constructor used to create the root of a routing tree. Root doesn't
// need any filters as it is invoked automatically by the server anyway.
// The routes will be added later, using Router's methods.
func New() *Router {
	return &Router{
		handler:    nil,
		fail:       DefaultFailHandler,
		routes:     nil,
		filters:    NewFilters(),
		middleware: make([]http.Handler, 0),
	}
}

// ServeHTTP method is here in order to ensure that Router implements the
// http.Handler interface. It is invoked automatically by http.Server if you
// assign Router in question as server's Handler. If this Router is not root,
// but a sub-router instead, its ServeHTTP method will be invoked by the parent
// Router whenever some request passes all its filters upon checkup.
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Cut path prefix (if set) from the reuqest URL path.
	if rtr.filters.PathPrefix != nil {
		r.URL.Path = strings.TrimPrefix(
			r.URL.Path, string(*rtr.filters.PathPrefix),
		)
	}

	// Parse path variables and alter http.Request.Context.
	r = rtr.vars(r)

	// Apply middleware.
	for _, mw := range rtr.middleware {
		mw.ServeHTTP(w, r)
	}

	// 1. Check if there are routes with matching filters.
	// 2. If not, use handler if present.
	// 3. If everything else failed, respond with a fail message.
	if sub, match := rtr.Match(r); match {
		sub.ServeHTTP(w, r)
	} else if rtr.handler != nil {
		rtr.handler.ServeHTTP(w, r)
	} else {
		rtr.fail.ServeHTTP(w, r)
	}
}

// Use registers a middleware handler on the Router.
func (rtr *Router) Use(h http.Handler) *Router {
	rtr.middleware = append(rtr.middleware, h)
	return rtr
}

// Use registers a middleware View handler on the Router.
func (rtr *Router) UseFunc(v View) *Router {
	rtr.middleware = append(rtr.middleware, v)
	return rtr
}

// Handler method sets router's handler.
func (rtr *Router) Handler(h http.Handler) *Router {
	rtr.handler = h
	return rtr
}

// HandleFunc method sets router's handler to a function.
func (rtr *Router) HandleFunc(v View) *Router {
	rtr.handler = v
	return rtr
}

// Fail method sets router's fail message.
func (rtr *Router) Fail(handler http.Handler) *Router {
	rtr.fail = handler
	return rtr
}

// FailFunc method sets router's fail message.
func (rtr *Router) FailFunc(v View) *Router {
	rtr.fail = v
	return rtr
}

// Subrouter method returns pointer to a new sub-router instance that inherits
// context from its parent.
func (rtr *Router) Subrouter() *Router {
	// Create new Router that inherits its parent's Context.
	sub := New()

	// Add it to parent's routes.
	rtr.routes = append(rtr.routes, sub)

	return sub
}

// Methods returns pointer to the same Router instance while altering its
// methods filter.
//
// NOTICE: If methods filter has already been set for this Router instance, it
// will get replaced!
func (rtr *Router) Methods(methods ...string) *Router {
	rtr.filters.Methods = NewMethodsFilter(methods...)
	return rtr
}

// Path returns pointer to the same Router instance while altering its path
// filter.
//
// NOTICE: This method replaces router's PathFilter with a newly created
// instance while setting PathPrefix to nil.
func (rtr *Router) Path(path string) *Router {
	rtr.filters.Path = NewPathFilter(path)
	rtr.filters.PathPrefix = nil
	return rtr
}

// PathPrefix returns pointer to the same Router instance while altering its
// path prefix filter.
//
// NOTICE: This method replaces router's PathPrefixFilter with a newly created
// instance while setting PathFilter to nil.
func (rtr *Router) PathPrefix(prefix string) *Router {
	rtr.filters.PathPrefix = NewPathPrefixFilter(prefix)
	rtr.filters.Path = nil
	return rtr
}

// Schemes returns pointer to the same Router instance while altering its
// schemes filter.
//
// NOTICE: This method replaces router's SchemesFilter with a newly created
// instance.
func (rtr *Router) Schemes(schemes ...string) *Router {
	for i, s := range schemes {
		schemes[i] = strings.ToLower(s)
	}
	rtr.filters.Schemes = NewSchemesFilter(schemes...)
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

		// Discarding all conversion errors in switch because we know
		// for sure that exp passed regex test for number.
		switch typ {
		case "int":
			vars[name], _ = strconv.Atoi(exp)

		case "nat":
			n, _ := strconv.ParseUint(exp, 10, 0)
			vars[name] = uint(n)

		case "str":
			vars[name] = exp

		default: // regex type
			vars[name] = exp
		}
	}

	return r.WithContext(context.WithValue(r.Context(), varsKey, vars))
}
