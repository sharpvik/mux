package mux

import (
	"net/http"
	"reflect"
)

// Filter is an interface type that represents essential functionality of a
// filter.
type Filter interface {
	Match(*http.Request) bool
}

// Filters is a concrete type that contains fields for every possible filter
// allowed on a Router. It ensures that only one filter of each type is used per
// Router instance.
type Filters struct {
	Methods *MethodsFilter // e.g. "GET", "POST", "PUT", "DELETE", etc.
	Path    *PathFilter    // e.g. "/public", "/api"
}

// NewFilters returns pointer to an empty set of filters.
func NewFilters() *Filters {
	return &Filters{nil, nil}
}

// Match method returns boolean value that tells you whether given request
// passed all the filters. Also, *Filters implements the Filter interface since
// it has this method.
func (fils *Filters) Match(r *http.Request) bool {
	v := reflect.ValueOf(*fils)

	// We'll have to go through every filter in the struct.
	for i := 0; i < v.NumField(); i++ {
		// Get reflect.Value of the i-th field in a struct.
		field := v.Field(i)

		// The nil filters are assumed to be all-permissive.
		if field.IsNil() {
			continue
		}

		// Type assertion to the Filter interface is needed.
		filter := field.Interface().(Filter)

		// Return false immediately if filter did not match the request.
		if !filter.Match(r) {
			return false
		}
	}

	// If all non-nil filters returned true, we return true.
	return true
}

// MethodsFilter takes care of filtering requests by method (e.g. "POST").
type MethodsFilter struct {
	// Methods contains a slice of all accepted methods. If you would like to
	// see all the ones that exist, go here:
	//
	//     https://golang.org/pkg/net/http/#pkg-constants
	//
	// It is advized that you use Go's standard "net/http" package in order to
	// manage these. For example:
	//
	//     package main
	//
	//     import (
	//         "net/http"
	//         "github.com/sharpvik/mux"
	//     )
	//
	//     func main() {
	//         // Create new filter instance.
	//         filter := mux.NewMethodsFilter(http.MethodPost)
	//
	//         // Add method "GET" to filter's Methods.
	//         filter.Methods = append(filter.Methods, http.MethodGet)
	//     }
	//
	Methods []string
}

// NewMethodsFilter function returns pointer to a custom MethodsFilter.
func NewMethodsFilter(methods ...string) *MethodsFilter {
	return &MethodsFilter{methods}
}

// Match method returns boolean value that tells you whether given request
// passed the filter. Also, *MethodsFilter implements the Filter interface since
// it has this method.
func (fil *MethodsFilter) Match(r *http.Request) bool {
	for _, m := range fil.Methods {
		if r.Method == m {
			return true
		}
	}
	return false
}

// PathFilter takes care of filtering requests by their URL path (e.g. "/api").
type PathFilter struct {
	// Path is a string that is used to decide whether given request matches
	// the URLFilter. It always begins with a /forward-slash.
	Path string
}

// NewPathFilter returns pointer to a newly created PathFilter.
func NewPathFilter(url string) *PathFilter {
	return &PathFilter{url}
}

// Match method returns boolean value that tells you whether given request
// passed the filter. Also, *PathFilter implements the Filter interface since
// it has this method.
func (fil *PathFilter) Match(r *http.Request) bool {
	return r.URL.Path == fil.Path
}
