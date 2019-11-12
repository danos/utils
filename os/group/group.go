// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// Specifically for the UnknownGroup(Id)Error sections:
//
// Copyright (c) 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MPL-2.0 and BSD-3-Clause

package group

import (
	"strconv"
)

//Group represents a group entry
type Group struct {
	Gid      uint32 // group id
	Name     string
	Password string
	Members  []string
}

// UnknownGroupIdError is returned by LookupId when
// a group cannot be found.
type UnknownGroupIdError int

func (e UnknownGroupIdError) Error() string {
	return "group: unknown groupid " + strconv.Itoa(int(e))
}

// UnknownGroupError is returned by Lookup when
// a group cannot be found.
type UnknownGroupError string

func (e UnknownGroupError) Error() string {
	return "group: unknown group " + string(e)
}
