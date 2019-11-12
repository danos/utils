// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2015 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package natsort

import (
	"testing"
)

var defaults = []string{"default2", "default", "default1"}
var thesame = []string{
	"foo1bar2",
	"foo1bar2",
	"foo1bar2",
	"foo1bar2",
	"foo1 bar2",
	"foo1 bar2",
	"foo1 bar2",
	"foo1 bar2",
}
var zs = []string{
	"z100.doc",
	"z1.doc",
	"z10.doc",
	"z102a.doc",
	"z101.doc",
	"z102b.doc",
	"z11.doc",
	"z12.doc",
	"z13.doc",
	"z14.doc",
	"z15.doc",
	"z16.doc",
	"z17.doc",
	"z18.doc",
	"z101.doc",
	"z19.doc",
	"z2.doc",
	"z20.doc",
	"z3.doc",
	"z4.doc",
	"z5.doc",
	"z6.doc",
	"z7.doc",
	"z8.doc",
	"z9.doc",
	"z.doc",
}

//list lifted from: http://www.davekoelle.com/alphanum.html
var others = []string{
	"1000X Radonius Maximus",
	"10X Radonius",
	"200X Radonius",
	"20X Radonius",
	"20X Radonius Prime",
	"30X Radonius",
	"40X Radonius",
	"Allegia 50 Clasteron",
	"Allegia 500 Clasteron",
	"Allegia 50B Clasteron",
	"Allegia 51 Clasteron",
	"Allegia 6R Clasteron",
	"Alpha 100",
	"Alpha 2",
	"Alpha 200",
	"Alpha 2A",
	"Alpha 2A-8000",
	"Alpha 2A-900",
	"Callisto Morphamax",
	"Callisto Morphamax 500",
	"Callisto Morphamax 5000",
	"Callisto Morphamax 600",
	"Callisto Morphamax 6000 SE",
	"Callisto Morphamax 6000 SE2",
	"Callisto Morphamax 700",
	"Callisto Morphamax 7000",
	"Xiph Xlater 10000",
	"Xiph Xlater 2000",
	"Xiph Xlater 300",
	"Xiph Xlater 40",
	"Xiph Xlater 5",
	"Xiph Xlater 50",
	"Xiph Xlater 500",
	"Xiph Xlater 5000",
	"Xiph Xlater 58",
}

var expectedDefaults = []string{
	"default",
	"default1",
	"default2",
}

var expectedZs = []string{
	"z.doc",
	"z1.doc",
	"z2.doc",
	"z3.doc",
	"z4.doc",
	"z5.doc",
	"z6.doc",
	"z7.doc",
	"z8.doc",
	"z9.doc",
	"z10.doc",
	"z11.doc",
	"z12.doc",
	"z13.doc",
	"z14.doc",
	"z15.doc",
	"z16.doc",
	"z17.doc",
	"z18.doc",
	"z19.doc",
	"z20.doc",
	"z100.doc",
	"z101.doc",
	"z101.doc",
	"z102a.doc",
	"z102b.doc",
}

var expectedOthers = []string{
	"10X Radonius",
	"20X Radonius",
	"20X Radonius Prime",
	"30X Radonius",
	"40X Radonius",
	"200X Radonius",
	"1000X Radonius Maximus",
	"Allegia 6R Clasteron",
	"Allegia 50 Clasteron",
	"Allegia 50B Clasteron",
	"Allegia 51 Clasteron",
	"Allegia 500 Clasteron",
	"Alpha 2",
	"Alpha 2A",
	"Alpha 2A-900",
	"Alpha 2A-8000",
	"Alpha 100",
	"Alpha 200",
	"Callisto Morphamax",
	"Callisto Morphamax 500",
	"Callisto Morphamax 600",
	"Callisto Morphamax 700",
	"Callisto Morphamax 5000",
	"Callisto Morphamax 6000 SE",
	"Callisto Morphamax 6000 SE2",
	"Callisto Morphamax 7000",
	"Xiph Xlater 5",
	"Xiph Xlater 40",
	"Xiph Xlater 50",
	"Xiph Xlater 58",
	"Xiph Xlater 300",
	"Xiph Xlater 500",
	"Xiph Xlater 2000",
	"Xiph Xlater 5000",
	"Xiph Xlater 10000",
}

func testLess(t *testing.T, a, b string) {
	if !Less(a, b) {
		t.Fatal(a, "should be less than", b)
	}
}

func testNotLess(t *testing.T, a, b string) {
	if Less(a, b) {
		t.Fatal(a, "should be less than", b)
	}
}

func TestDefaultString(t *testing.T) {
	a := "default"
	b := "default2"
	testLess(t, a, b)
	testNotLess(t, b, a)
}

func TestDifferentNumbers(t *testing.T) {
	a := "default1"
	b := "default2"
	testLess(t, a, b)
	testNotLess(t, b, a)
}

func TestTheSameString(t *testing.T) {
	a := "default"
	b := "default"
	testLess(t, a, b)
	testLess(t, b, a)
}

func TestMultiNumber(t *testing.T) {
	a := "default1a2"
	b := "default1a3"
	testLess(t, a, b)
	testNotLess(t, b, a)
}

func eqLists(l1, l2 []string) bool {
	if len(l2) != len(l1) {
		return false
	}

	for i, elem := range l1 {
		if l2[i] != elem {
			return false
		}
	}
	return true
}

func testList(t *testing.T, in, expected []string) {
	Sort(in)
	if !eqLists(in, expected) {
		t.Fatal("lists did not match")
	}
}

func TestDefaults(t *testing.T) {
	testList(t, defaults, expectedDefaults)
}

func TestZs(t *testing.T) {
	testList(t, zs, expectedZs)
}
func TestOthers(t *testing.T) {
	testList(t, others, expectedOthers)
}
func TestTheSame(t *testing.T) {
	testList(t, thesame, thesame)
}
