// Copyright (c) 2019, AT&T Intellectual Property.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0
package args

import "fmt"

func ExampleParseArgs_quoted() {
	fmt.Printf("%q\n", ParseArgs("foo \"foo bar\" baz"))
	// Output: ["foo" "foo bar" "baz"]
}

func ExampleParseArgs_mixedQuotes() {
	fmt.Printf("%q\n", ParseArgs("foo \"foo'bar\" baz"))
	// Output: ["foo" "foo'bar" "baz"]
}

func ExampleParseArgs_escaped() {
	fmt.Printf("%q\n", ParseArgs(`foo "foo'bar" \ baz`))
	// Output: ["foo" "foo'bar" " baz"]
}

func ExampleParseArgs_quoteInArg() {
	fmt.Printf("%q\n", ParseArgs(`fo"o"o "foo'bar" \ baz`))
	// Output: ["fooo" "foo'bar" " baz"]
}
