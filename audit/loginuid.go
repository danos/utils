package audit

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func GetPidLoginuid(pid int32) (uint32, error) {
	procfile := fmt.Sprintf("/proc/%d/loginuid", pid)
	f, e := os.Open(procfile)
	if e != nil {
		return 0, e
	}
	defer f.Close()

	var buf []byte = make([]byte, 10)
	n, e := f.Read(buf)
	if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
		return 0, e
	}

	uids := string(buf[:n])
	uid, e := strconv.Atoi(uids)
	if e != nil {
		return 0, e
	}

	return uint32(uid), nil
}
