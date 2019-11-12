// Copyright (c) 2019, AT&T Intellectual Property.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package tty

import (
	"fmt"
	"os"
	"strings"
)

func TtyNameForPid(pid int) (string, error) {
	tty, err := os.Readlink(fmt.Sprintf("/proc/%v/fd/0", pid))
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(tty, "/dev/"), nil
}
