// Copyright (c) 2018-2019, AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// +build linux
// +build cgo

package audit

/*
#cgo LDFLAGS: -laudit
#include <stdlib.h>
#include <libaudit.h>

static void logUserCmd(char *msg, uint result) {
	int audit_fd = audit_open();
	if (audit_fd < 0) {
		return;
	}
	audit_log_user_message(audit_fd, AUDIT_USER_CMD, msg, NULL, NULL, NULL, result);
	audit_close(audit_fd);
}

static void logUserConfig(char *msg, uint result) {
	int audit_fd = audit_open();
	if (audit_fd < 0) {
		return;
	}
	audit_log_user_message(audit_fd, AUDIT_USYS_CONFIG, msg, NULL, NULL, NULL, result);
	audit_close(audit_fd);
}
*/
import "C"
import "unsafe"

func LogUserCmd(msg string, result uint) {
	cstr := C.CString(msg)
	C.logUserCmd(cstr, C.uint(result))
	C.free(unsafe.Pointer(cstr))
}

func LogUserConfig(msg string, result bool) {
	var res uint
	if result {
		res = 1
	} else {
		res = 0
	}
	cstr := C.CString(msg)
	C.logUserConfig(cstr, C.uint(res))
	C.free(unsafe.Pointer(cstr))
}
