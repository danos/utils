// Copyright (c) 2020, AT&T Intellectual Property. All rights reserved.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// Package shell provides helpers for working with shell scripts
package shell

import "strings"

const (
	specialChars    = "#\\&><!|*}{)(][:?^;"
	whitespaceChars = "\011\012\013\014\015 "
	strongChars     = "$\""
	weakChars       = "'"
)

// Quote will properly quote strings such that the shell will not interpret them.
func Quote(in string) string {
	hasSpecial := strings.ContainsAny(in,
		strongChars+weakChars+specialChars+whitespaceChars)
	needsStrongQuote := strings.ContainsAny(in, strongChars)
	needsWeakQuote := strings.Contains(in, weakChars)
	switch {
	case needsStrongQuote && needsWeakQuote:
		out := strings.Replace(in, "'", "\\'", -1)
		return "$'" + out + "'"
	case needsStrongQuote:
		return "'" + in + "'"
	case needsWeakQuote:
		return "\"" + in + "\""
	case hasSpecial:
		return "'" + in + "'"
	default:
		return in
	}
}
