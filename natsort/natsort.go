// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package natsort

import (
	"sort"
	"strconv"
	"unicode"
)

type natsorter []string

func (s natsorter) Len() int           { return len(s) }
func (s natsorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s natsorter) Less(i, j int) bool { return Less(s[i], s[j]) }

func split(s string) []string {
	out := make([]string, 0, 3)
	var indigit bool
	var start, pos int
	var r rune
	for pos, r = range s {
		if unicode.IsDigit(r) {
			if pos > start && !indigit {
				out = append(out, s[start:pos])
				start = pos
			}
			indigit = true
		} else {
			if pos > start && indigit {
				out = append(out, s[start:pos])
				start = pos
			}
			indigit = false
		}
	}
	out = append(out, s[start:])
	return out
}

func Less(ain, bin string) (out bool) {
	if ain == bin {
		return true
	}
	acomp := split(ain)
	bcomp := split(bin)
	for i, a := range acomp {
		if i >= len(bcomp) {
			return false
		}
		b := bcomp[i]
		if a == b {
			continue
		}
		if aint, err := strconv.Atoi(a); err == nil {
			if bint, err := strconv.Atoi(b); err == nil {
				return aint < bint
			}
		}
		return ain < bin
	}
	return true
}

func Sort(sl []string) {
	sort.Sort(natsorter(sl))
}
