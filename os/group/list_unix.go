// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// +build darwin freebsd linux netbsd openbsd
// +build cgo

package group

import (
	"runtime"
	"sync"
)

/*
#include <unistd.h>
#include <grp.h>
#include <stdlib.h>
*/
import "C"

var grmu = sync.Mutex{}

func list() ([]*Group, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	grmu.Lock()
	defer grmu.Unlock()
	var groups = make([]*Group, 0)
	var result *C.struct_group
	C.setgrent()
	for {
		result = C.getgrent()
		if result == nil {
			break
		}
		g := newGroup(result)
		groups = append(groups, g)
	}
	C.endgrent()
	return groups, nil
}
