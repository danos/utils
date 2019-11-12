// Copyright (c) 2018-2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014-2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package pathutil

import (
	"fmt"
	"net/url"
	"strings"
)

type PathElementAttrs struct {
	Secret bool
}

func NewPathElementAttrs() PathElementAttrs {
	var attrs PathElementAttrs
	attrs.Secret = false
	return attrs
}

type PathAttrs struct {
	Attrs []PathElementAttrs
}

func NewPathAttrs() PathAttrs {
	var attrs PathAttrs
	attrs.Attrs = make([]PathElementAttrs, 0)
	return attrs
}

func Makepath(path_string string) []string {
	raw_path := strings.Split(path_string, "/")
	path := make([]string, 0, len(raw_path))

	for _, v := range raw_path {
		if v != "" {
			new_v, _ := url.QueryUnescape(v)
			path = append(path, new_v)
		}
	}
	return path
}

func Pathstr(path []string) string {
	var str string
	for _, v := range path {
		//Make this compatible with the liburiparser c library
		//liburiparser does not handle '+' as spaces as specified in RFC-3986
		str += "/" + strings.Replace(url.QueryEscape(v), "+", "%20", -1)
	}
	return str
}

func Copypath(path []string) []string {
	npath := make([]string, len(path))
	copy(npath, path)
	return npath
}

func CopyAppend(path []string, vals ...string) []string {
	npath := make([]string, len(path), len(path)+len(vals))
	copy(npath, path)
	return append(npath, vals...)
}

func RedactPath(path []string, pathAttrs *PathAttrs) ([]string, error) {
	if pathAttrs == nil || len(path) != len(pathAttrs.Attrs) {
		return []string{"<path redaction failed>"}, fmt.Errorf("Unable to redact command")
	}

	/* Redact values for all secret elements in the path */
	rpath := make([]string, len(path))
	for i, v := range pathAttrs.Attrs {
		if v.Secret {
			rpath[i] = "**"
		} else {
			rpath[i] = path[i]
		}
	}

	return rpath, nil
}
