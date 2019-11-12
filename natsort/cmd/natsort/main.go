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
	"github.com/danos/utils/natsort"
	"os"
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

	scanner := bufio.NewScanner(bufio.NewReader(f))
	scanner.Split(bufio.ScanWords)
	elems := make([]string, 0)

	for scanner.Scan() {
		t := scanner.Text()
		elems = append(elems, t)
	}

	natsort.Sort(elems)
	for _, e := range elems {
		fmt.Println(e)
	}
}
