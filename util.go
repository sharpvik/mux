package mux

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Vars function returns path variables in a map[string]interface{} and a
// boolean success confirmation flag.
func Vars(r *http.Request) (varsmap map[string]interface{}, ok bool) {
	v := r.Context().Value(varsKey)
	if ok = v != nil; ok {
		varsmap = v.(map[string]interface{})
		return
	}
	return
}

// isVar tells you whether this path segment pattern was intended as a variable.
// The pattern is either an arbitrary string or of "{varname:vartype}" form.
func isVar(pattern string) bool {
	regex := regexp.MustCompile(`\{\w+:(int|str)\}`)
	return regex.MatchString(pattern)
}

// varData returns path var's name and type from given pattern where pattern is
// something like "{id:int}".
func varData(pattern string) (name string, typ pathVarType) {
	trim := string([]rune(pattern)[1:strings.IndexRune(pattern, '}')])
	split := strings.Split(trim, ":")
	name = split[0]
	typeString := split[1]

	switch typeString {
	case "int":
		typ = pint
	case "str":
		typ = pstr
	default:
		panic(fmt.Sprintf("invalid type in path %s", pattern))
	}

	return
}
