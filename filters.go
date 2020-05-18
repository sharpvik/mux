package mux

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// Filter is an interface type that represents functionality of a filter.
type Filter interface {
	Match(*http.Request) bool
}

// Filters is a concrete type that contains fields for every possible filter
// allowed on a Router. It ensures that only one filter of each type is used per
// Router instance.
type Filters struct {
	Schemes    *SchemesFilter    // e.g. "http" or "https".
	Methods    *MethodsFilter    // e.g. "GET", "POST", "PUT", "DELETE", etc.
	Path       *PathFilter       // e.g. "/home" or "/r/{sub:str}/{id:int}".
	PathPrefix *PathPrefixFilter // e.g. "/api".
}

// NewFilters returns pointer to an empty set of filters.
func NewFilters() *Filters {
	return &Filters{nil, nil, nil, nil}
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
// If you would like to see all the request methods that exist, go here:
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
//         // Create new filter instance using a constant from "net/http".
//         filter := mux.NewMethodsFilter(http.MethodPost)
//
//         // Add method "GET" to filter's Methods.
//         filter.Methods.Add(http.MethodGet)
//     }
//
type MethodsFilter struct {
	Methods set
}

// NewMethodsFilter function returns pointer to a custom MethodsFilter.
func NewMethodsFilter(methods ...string) *MethodsFilter {
	return &MethodsFilter{newSet(methods...)}
}

// Match method returns boolean value that tells you whether given request
// passed the filter. Also, *MethodsFilter implements the Filter interface since
// it has this method.
func (fil MethodsFilter) Match(r *http.Request) bool {
	return fil.Methods.Has(r.Method)
}

// PathFilter takes care of filtering requests by their URL path (e.g. "/api").
type PathFilter struct {
	// Path is a pattern string that is used to compose and compile a proper
	// regual expression (Regexp) that will be used to match URL paths.
	Path string

	// Regexp is a compiled regular expression that is created by the
	// NewPathFilter function; it is going to be used to check if request path
	// matches the PathFilter.
	Regexp *regexp.Regexp

	// hasVars is a boolean flag that tells us whether this PathFilter had path
	// variables in its template path.
	hasVars bool
}

// NewPathFilter returns pointer to a newly created PathFilter. It also ensures
// that the first character in the uri is a forward-slash -- if it isn't there,
// it will be inserted.
func NewPathFilter(path string) *PathFilter {
	// Create a dummy PathFilter.
	fil := &PathFilter{"", nil, false}

	// Ensure that the leading slash is present in the path.
	if []byte(path)[0] != '/' {
		path = "/" + path
	}
	fil.Path = path

	// Split path template by "/" and build an appropriate regular expression.
	split := strings.Split(path, "/")[1:]
	var exp string

	for _, e := range split {
		if isVar(e) {
			fil.hasVars = true

			_, typ := varData(e)
			sub := "/"
			switch typ {
			case "int":
				sub = sub + `(-?[1-9]\d*|0)`

			case "str":
				sub = sub + `[a-zA-Z_]+`

			case "nat":
				sub = sub + `([1-9]\d*|0)`

			default: // regex type
				sub = sub + typ
			}

			exp = exp + sub
		} else {
			exp = exp + "/" + e
		}
	}

	// Try to compile generated regular expression. Panic if that fails.
	regex, err := regexp.Compile(exp)
	if err != nil {
		panic(fmt.Sprintf("can't compile regex %s: %v", exp, err))
	}
	fil.Regexp = regex

	return fil
}

// Match method returns boolean value that tells you whether given request
// passed the filter. Also, *PathFilter implements the Filter interface since
// it has this method.
func (fil *PathFilter) Match(r *http.Request) bool {
	return fil.Regexp.MatchString(r.URL.Path)
}

// PathPrefixFilter takes care of filtering requests by URL path prefix.
// It is an alias to the standard string type. The string it wraps is the
// aforementioned path prefix which we wish to utilize for route matching
// purposes.
type PathPrefixFilter string

// NewPathPrefixFilter returns reference to a newly created PathPrefixFilter.
func NewPathPrefixFilter(prefix string) *PathPrefixFilter {
	fil := PathPrefixFilter(prefix)
	return &fil
}

// Match method uses the string (that PathPrefixFilter wraps around) to decide
// whether the request in question matches or not.
func (fil *PathPrefixFilter) Match(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, string(*fil))
}

// SchemesFilter takes care of filtering requests by scheme (e.g. "https").
type SchemesFilter struct {
	Schemes set
}

// NewSchemesFilter function returns pointer to a custom SchemesFilter.
func NewSchemesFilter(schemes ...string) *SchemesFilter {
	return &SchemesFilter{newSet(schemes...)}
}

// Match method returns boolean value that tells you whether given request
// passed the filter. Also, *SchemesFilter implements the Filter interface since
// it has this method.
func (fil *SchemesFilter) Match(r *http.Request) bool {
	scheme := r.URL.Scheme

	if scheme == "" {
		if r.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
	}

	return fil.Schemes.Has(scheme)
}
