package mux

import (
	"fmt"
	"regexp"
	"strings"
)

func isVar(pattern string) bool {
	regex := regexp.MustCompile(`\{\w+:(int|str)\}`)
	return regex.MatchString(pattern)
}

func varData(pattern string) (name string, typ pathVarType) {
	trim := string([]rune(pattern)[1:strings.IndexRune(pattern, '}')])
	split := strings.Split(trim, ":")
	name = split[0]
	typeString := split[1]

	switch typeString {
	case "int":
		typ = Int
	case "str":
		typ = Str
	default:
		panic(fmt.Sprintf("invalid type in path %s", pattern))
	}

	return
}
