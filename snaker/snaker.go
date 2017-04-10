// Package snaker provides methods to convert CamelCase to and from snake_case.
//
// snaker takes into takes into consideration common initialisms (ie, ID, HTTP,
// ACL, etc) when converting to/from CamelCase and snake_case.
package snaker

//go:generate ./gen.sh --update

import (
	"fmt"
	"strings"
	"unicode"
)

// CamelToSnake converts s to snake_case.
func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	rs := []rune(s)

	var r string
	var lastWasUpper, lastWasLetter, lastWasIsm, isUpper, isLetter bool
	for i := 0; i < len(rs); {
		isUpper = unicode.IsUpper(rs[i])
		isLetter = unicode.IsLetter(rs[i])

		// append _ when last was not upper and not letter
		if (lastWasLetter && isUpper) || (lastWasIsm && isLetter) {
			r += "_"
		}

		// determine next to append to r
		var next string
		if ism := peekInitialism(rs[i:]); ism != "" && (!lastWasUpper || lastWasIsm) {
			next = ism
		} else {
			next = string(rs[i])
		}

		// save for next iteration
		lastWasIsm = false
		if len(next) > 1 {
			lastWasIsm = true
		}
		lastWasUpper = isUpper
		lastWasLetter = isLetter

		r += next
		i += len(next)
	}

	return strings.ToLower(r)
}

// CamelToSnakeIdentifier converts s to its snake_case identifier.
func CamelToSnakeIdentifier(s string) string {
	return toIdentifier(CamelToSnake(s))
}

// SnakeToCamel converts s to CamelCase.
func SnakeToCamel(s string) string {
	var r string

    if len(s) == 0{
        return s
    }

    //ME: hack snake just for those of having "_"
    if strings.Index(s,"_") < 0 {
        return strings.ToUpper(s[:1]) + s[1:]
    }

	for _, w := range strings.Split(s, "_") {
		if w == "" {
			continue
		}

		u := strings.ToUpper(w)
		if ok := commonInitialisms[u]; ok {
			r += u
		} else {
			r += strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}

	return r
}

// SnakeToCamelIdentifier converts s to its CamelCase identifier (first
// letter is capitalized).
func SnakeToCamelIdentifier(s string) string {
	return SnakeToCamel(toIdentifier(s))
}

// AddInitialisms adds initialisms to the recognized initialisms.
func AddInitialisms(initialisms ...string) error {
	for _, s := range initialisms {
		if len(s) < minInitialismLen || len(s) > maxInitialismLen {
			return fmt.Errorf("%s does not have length between %d and %d", s, minInitialismLen, maxInitialismLen)
		}
		commonInitialisms[s] = true
	}

	return nil
}
