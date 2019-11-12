// Copyright (c) 2019, AT&T Intellectual Property.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package args

import "unicode"

// ParseArgs takes a string and returns a set of arguments preserving
// quoted values as a single argument. This is similar to shell argument
// parsing.
func ParseArgs(s string) []string {
	const (
		inArg = iota
		inQuote
		inWhitespace
	)
	state := inWhitespace
	quoteChar := rune(0)
	arg := ""
	argv := []string{}

	isQuote := func(c rune) bool {
		return c == '"' || c == '\''
	}

	isEscape := func(c rune) bool {
		return c == '\\'
	}

	isWhitespace := func(c rune) bool {
		return unicode.IsSpace(c)
	}

	runes := []rune(s)
	strLen := len(runes)
	var c rune
	for i := 0; i < strLen; i++ {
		c = runes[i]
		switch {
		case isQuote(c):
			switch state {
			case inWhitespace:
				arg = ""
				state = inQuote
				quoteChar = c
			case inArg:
				state = inQuote
				quoteChar = c
			case inQuote:
				if c == quoteChar {
					state = inArg
				} else {
					arg = arg + string(c)
				}
			}
		case isWhitespace(c):
			switch state {
			case inArg:
				argv = append(argv, arg)
				state = inWhitespace
			case inQuote:
				arg = arg + string(c)
			case inWhitespace:
			}
		case isEscape(c):
			if state == inWhitespace {
				arg = ""
				state = inArg
			}
			if i == strLen-1 {
				panic("unexpected end of string")
			} else {
				i++
				c = runes[i]
				arg = arg + string(c)
			}
		default:
			switch state {
			case inArg, inQuote:
				arg = arg + string(c)
			case inWhitespace:
				arg = ""
				arg = arg + string(c)
				state = inArg
			}
		}
	}

	switch state {
	case inArg:
		argv = append(argv, arg)
	case inQuote:
		panic("unexpected end of string")
	}
	return argv
}
