// Copyright (c) 2019, AT&T Intellectual Property.
// All rights reserved.
//
// Copyright (c) 2014-2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package patherr

import (
	"bytes"
	"fmt"
	"github.com/danos/utils/natsort"
	"strings"
	"text/tabwriter"
)

const cfgPath = `Configuration path: `
const operCommand = `Command: `

type CommandInval struct {
	Path []string
	Fail string
}

func (e *CommandInval) Error() string {
	if len(e.Path) == 0 {
		return fmt.Sprintf("Invalid command: [%s]", e.Fail)
	}
	return fmt.Sprintf("Invalid command: %s [%s]", strings.Join(e.Path, " "), e.Fail)
}

type PathInval struct {
	Path        []string
	Fail        string
	Operational bool
}

func (e *PathInval) Error() string {
	prefix := cfgPath
	if e.Operational {
		prefix = operCommand
	}
	if len(e.Path) == 0 {
		return fmt.Sprintf("%s [%s] is not valid", prefix, e.Fail)
	}
	return fmt.Sprintf("%s %s [%s] is not valid", prefix, strings.Join(e.Path, " "), e.Fail)
}

type PathAmbig struct {
	Path        []string
	Fail        string
	Matches     map[string]string
	Operational bool
}

func (e *PathAmbig) Error() string {
	var buf = new(bytes.Buffer)
	twriter := tabwriter.NewWriter(buf, 8, 0, 1, ' ', 0)

	prefix := cfgPath
	if e.Operational {
		prefix = operCommand
	}

	if len(e.Path) == 0 {
		fmt.Fprintf(buf, "%s [%s] is ambiguous\n", prefix, e.Fail)
	} else {
		fmt.Fprintf(buf, "%s %s [%s] is ambiguous\n", prefix, strings.Join(e.Path, " "), e.Fail)
	}
	fmt.Fprintf(buf, "\n  Possible completions:\n")

	sorted := make([]string, 0, len(e.Matches))
	for n, _ := range e.Matches {
		sorted = append(sorted, n)
	}

	natsort.Sort(sorted)
	for i, name := range sorted {
		fmt.Fprintf(twriter, "    %s\t%s", name, e.Matches[name])
		if i != len(sorted)-1 {
			fmt.Fprintf(twriter, "\n")
		}
	}
	twriter.Flush()
	return buf.String()
}
