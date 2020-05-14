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
	return regexp.MustCompile(`\{\w+:.+\}`).MatchString(pattern)
}

// varData returns path var's name and type from given pattern where pattern is
// something like "{id:int}".
func varData(pattern string) (name string, typ string) {
	trim := string([]rune(pattern)[1 : len(pattern)-1])
	split := strings.Split(trim, ":")
	name = split[0]
	typ = split[1]

	switch typ {
	case "int", "str", "nat": // NOP case just to catch regex in typ.
	default:
		// At this point we assume that it's either a regex expression that can
		// be compiled, or an invalid type (in which case we should panic).
		_, err := regexp.Compile(typ)
		if err != nil {
			panic(fmt.Sprintf("invalid type/regex in path %s", pattern))
		}
	}

	return
}
