// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014, 2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// +build darwin freebsd linux netbsd openbsd
// +build cgo

package group

import (
	"os/user"
	"strconv"
	"unsafe"
)

/*
#include <grp.h>
#include <stdlib.h>

static char *getMember(char **gr_mem, int n) {
	return gr_mem[n];
}

static int my_getgrouplist(const char *user, int group, int *groups,
				int *ngroups, int *size) {
	if (getgrouplist(user, group, groups, ngroups) == -1) {
		// FIXME: this is only because of missing C-sizeof -dgollub
		*size = sizeof(gid_t) * (*ngroups);
		return -1;
	}
	return 0;
}

static int getGid(int *gids, int n) {
	return gids[n];
}
*/
import "C"

func lookupUser(username string) ([]*Group, error) {
	var ngroups C.int = 0
	var size C.int
	var groups = make([]*Group, 0)
	var cusername *C.char = C.CString(username)
	defer C.free(unsafe.Pointer(cusername))

	// Lookup user for getgroups which requires primary GID
	user, e := user.Lookup(username)
	if e != nil {
		return nil, e
	}

	// Dummy space to determine the actual number of groups (ngroups)
	dummy_groups := C.malloc(C.size_t(0))
	defer C.free(dummy_groups)

	my_gid, e := strconv.Atoi(user.Gid)
	if e != nil {
		return nil, e
	}

	// Determine the number of groups the user belongs to
	C.my_getgrouplist(cusername, C.int(my_gid), (*C.int)(dummy_groups),
		&ngroups, &size)

	// Allocate the required space.
	// XXX: Since I could not quickly enough find sizeof(gid_t) for Go
	// I introduced to the getgrouplist wrapper this additional size argument
	// to allocate the correct amount of memory -dgollub
	my_groups := C.malloc(C.size_t(size))
	defer C.free(my_groups)

	C.my_getgrouplist(cusername, C.int(my_gid), (*C.int)(my_groups),
		&ngroups, &size)

	// Lookup the group names based on the GIDs and assemble a list
	// of group names
	for i := 0; i < int(ngroups); i++ {
		gid := C.getGid((*C.int)(my_groups), C.int(i))

		g, e := lookupUnixGid(int(gid))
		if e != nil {
			return nil, e
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func lookupUid(uid string) ([]*Group, error) {
	u, e := user.LookupId(uid)
	if e != nil {
		return nil, e
	}

	return lookupUser(u.Username)
}

func lookup(groupname string) (*Group, error) {
	return lookupGroup(groupname)
}

func lookupId(gid string) (*Group, error) {
	i, e := strconv.Atoi(gid)
	if e != nil {
		return nil, e
	}
	return lookupUnixGid(i)
}

func newGroup(grp *C.struct_group) *Group {
	g := &Group{
		Gid:      uint32(grp.gr_gid),
		Name:     C.GoString(grp.gr_name),
		Password: C.GoString(grp.gr_passwd),
		Members:  make([]string, 0),
	}
	for i := 0; ; i++ {
		mem := C.getMember(grp.gr_mem, C.int(i))
		if mem == nil {
			break
		}
		g.Members = append(g.Members, C.GoString(mem))
	}
	return g
}
