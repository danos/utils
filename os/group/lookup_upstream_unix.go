// Copyright 2011 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// SPDX-License-Identifier: BSD-3-Clause

// +build darwin freebsd linux netbsd openbsd
// +build cgo

package group

import (
	"fmt"
	"syscall"
	"unsafe"
)

/*
#include <unistd.h>
#include <grp.h>
#include <stdlib.h>

static int mygetgrgid_r(int gid, struct group *grp,
	char *buf, size_t buflen, struct group **result) {
	return getgrgid_r(gid, grp, buf, buflen, result);
}

static int mygetgrnam_r(const char *name, struct group *grp,
	char *buf, size_t buflen, struct group **result) {
	return getgrnam_r(name, grp, buf, buflen, result);
}
*/
import "C"

func lookupGroup(groupname string) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group

	buf := alloc(groupBuffer)
	defer buf.free()
	cname := C.CString(groupname)
	defer C.free(unsafe.Pointer(cname))

	err := retryWithBuffer(buf, func() syscall.Errno {
		return syscall.Errno(C.mygetgrnam_r(cname,
			&grp,
			(*C.char)(buf.ptr),
			C.size_t(buf.size),
			&result))
	})
	if err != nil {
		return nil, fmt.Errorf("user: lookup groupname %s: %v", groupname, err)
	}
	if result == nil {
		return nil, UnknownGroupError(groupname)
	}
	return newGroup(&grp), nil
}

func lookupUnixGid(gid int) (*Group, error) {
	var grp C.struct_group
	var result *C.struct_group

	buf := alloc(groupBuffer)
	defer buf.free()

	err := retryWithBuffer(buf, func() syscall.Errno {
		// mygetgrgid_r is a wrapper around getgrgid_r to
		// to avoid using gid_t because C.gid_t(gid) for
		// unknown reasons doesn't work on linux.
		return syscall.Errno(C.mygetgrgid_r(C.int(gid),
			&grp,
			(*C.char)(buf.ptr),
			C.size_t(buf.size),
			&result))
	})
	if err != nil {
		return nil, fmt.Errorf("user: lookup groupid %d: %v", gid, err)
	}
	if result == nil {
		return nil, UnknownGroupIdError(gid)
	}
	return newGroup(&grp), nil
}

type bufferKind C.int

const (
	groupBuffer = bufferKind(C._SC_GETGR_R_SIZE_MAX)
)

func (k bufferKind) initialSize() C.size_t {
	sz := C.sysconf(C.int(k))
	if sz == -1 {
		// DragonFly and FreeBSD do not have _SC_GETPW_R_SIZE_MAX.
		// Additionally, not all Linux systems have it, either. For
		// example, the musl libc returns -1.
		return 1024
	}
	if !isSizeReasonable(int64(sz)) {
		// Truncate.  If this truly isn't enough, retryWithBuffer will error on the first run.
		return maxBufferSize
	}
	return C.size_t(sz)
}

type memBuffer struct {
	ptr  unsafe.Pointer
	size C.size_t
}

func alloc(kind bufferKind) *memBuffer {
	sz := kind.initialSize()
	return &memBuffer{
		ptr:  C.malloc(sz),
		size: sz,
	}
}

func (mb *memBuffer) resize(newSize C.size_t) {
	mb.ptr = C.realloc(mb.ptr, newSize)
	mb.size = newSize
}

func (mb *memBuffer) free() {
	C.free(mb.ptr)
}

// retryWithBuffer repeatedly calls f(), increasing the size of the
// buffer each time, until f succeeds, fails with a non-ERANGE error,
// or the buffer exceeds a reasonable limit.
func retryWithBuffer(buf *memBuffer, f func() syscall.Errno) error {
	for {
		errno := f()
		if errno == 0 {
			return nil
		} else if errno != syscall.ERANGE {
			return errno
		}
		newSize := buf.size * 2
		if !isSizeReasonable(int64(newSize)) {
			return fmt.Errorf("internal buffer exceeds %d bytes", maxBufferSize)
		}
		buf.resize(newSize)
	}
}

const maxBufferSize = 1 << 20

func isSizeReasonable(sz int64) bool {
	return sz > 0 && sz <= maxBufferSize
}
