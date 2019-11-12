// Copyright (c) 2017-2019, AT&T Intellectual Property.
// All rights reserved.
// Copyright (c) 2014-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"fmt"
	"os"
	"strings"

	"github.com/danos/utils/pathutil"
	"os/exec"
)

var NewExecError func(path []string, err string) error

func init() {
	NewExecError = newExecError
}

type execError struct {
	path []string
	err  string
}

func (e *execError) Error() string {
	return fmt.Sprintf("%s\n%s\n", e.path, e.err)
}

func newExecError(path []string, err string) error {
	return &execError{path: path, err: err}
}

type Output struct {
	Path   []string
	Output string
}

func (e *Output) String() string {
	return fmt.Sprintf("%s\n%s\n", e.Path, e.Output)
}

type ExecFunc func() ([]*Output, []error, bool)

func AppendOutput(fn ExecFunc, outs []*Output, errs []error) ([]*Output, []error, bool) {
	couts, cerrs, ok := fn()
	outs = append(outs, couts...)
	errs = append(errs, cerrs...)
	return outs, errs, ok
}

func ExecNoErr(env, path []string, cmd string) (*Output, error) {
	fn := (*exec.Cmd).Output
	return execfn(fn, env, path, cmd)
}
func Exec(env, path []string, cmd string) (*Output, error) {
	fn := (*exec.Cmd).CombinedOutput
	return execfn(fn, env, path, cmd)
}
func execfn(fn func(*exec.Cmd) ([]byte, error), env, path []string, cmd string) (*Output, error) {
	var args []string
	var c *exec.Cmd
	path = pathutil.Copypath(path)
	cliexec := "/opt/vyatta/bin/cliexec"
	if _, err := os.Stat(cliexec); os.IsNotExist(err) {
		args = strings.Split(cmd, " ")
		c = exec.Command(args[0], args[1:]...)
	} else {
		args = []string{"-s", cmd}
		c = exec.Command(cliexec, args...)
	}

	c.Env = env

	msg, e := fn(c)
	if e != nil {
		if _, ok := e.(*exec.ExitError); !ok {
			return nil, NewExecError(path, e.Error())
		}
	}
	if c.ProcessState.Success() {
		if len(msg) == 0 {
			return nil, nil
		}
		return &Output{Path: path, Output: string(msg)}, nil
	}

	return nil, NewExecError(path, string(msg))
}

func EnvWithSocket(sid string, path []string, action string, cact string, socket string) []string {
	var env = []string{
		"vyatta_htmldir=/opt/vyatta/share/html",
		"vyatta_datadir=/opt/vyatta/share",
		"vyatta_op_templates=/opt/vyatta/share/vyatta-op/templates",
		"vyatta_sysconfdir=/opt/vyatta/etc",
		"vyatta_sharedstatedir=/opt/vyatta/com",
		"vyatta_sbindir=/opt/vyatta/sbin",
		"vyatta_cfg_templates=/opt/vyatta/share/vyatta-cfg/templates",
		"vyatta_bindir=/opt/vyatta/bin",
		"vyatta_libdir=/opt/vyatta/lib",
		"vyatta_localstatedir=/opt/vyatta/var",
		"vyatta_libexecdir=/opt/vyatta/libexec",
		"vyatta_prefix=/opt/vyatta",
		"vyatta_datarootdir=/opt/vyatta/share",
		"vyatta_configdir=/opt/vyatta/config",
		"vyatta_infodir=/opt/vyatta/share/info",
		"vyatta_localedir=/opt/vyatta/share/locale",
		"PATH=/usr/local/bin:/usr/bin:/bin:/usr/local/sbin:/usr/sbin:/sbin:/opt/vyatta/bin:/opt/vyatta/bin/sudo-users:/opt/vyatta/sbin",
		"PERL5LIB=/opt/vyatta/share/perl5",
	}

	if sid != "" {
		env = append(env, "VYATTA_CONFIG_SID="+sid)
	}
	if cact != "" {
		env = append(env, "COMMIT_ACTION="+cact)
	}
	if socket != "" {
		env = append(env, "VYATTA_CONFIG_SOCKET="+socket)
	}
	env = append(env, "CONFIGD_PATH="+pathutil.Pathstr(path))
	env = append(env, "CONFIGD_EXT="+action)
	return env
}

func Env(sid string, path []string, action string, cact string) []string {
	return EnvWithSocket(sid, path, action, cact, "")
}
