// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/danos/utils/tsort"
)

func main() {
	var e error

	f := os.Stdin
	n := os.Args[0]
	if len(os.Args) > 1 {
		fname := os.Args[1]
		f, e = os.Open(fname)
		if e != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", n, e)
			os.Exit(1)
		}
	}
	g := tsort.New()

	scanner := bufio.NewScanner(bufio.NewReader(f))
	scanner.Split(bufio.ScanWords)

	var v, w string
	elems := make([]string, 0)

	for scanner.Scan() {
		t := scanner.Text()
		elems = append(elems, t)
	}
	if len(elems)%2 != 0 {
		fmt.Fprintf(os.Stderr, "%s: odd data count\n", n)
		os.Exit(1)
	}
	for i := 0; i < len(elems); i++ {
		v = elems[i]
		i += 1
		w = elems[i]
		g.AddEdge(v, w)
	}
	out, err := g.Sort()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, s := range out {
		fmt.Println(s)
	}
}
