// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"reflect"
	"sort"
	"testing"
)

type LogType uint

const (
	LOG_TYPE_USER_CFG = iota
	LOG_TYPE_USER_CMD = iota
)

type UserLog struct {
	Type   LogType
	Msg    string
	Result uint
}

type UserLogSlice []UserLog

func (s UserLogSlice) Len() int {
	return len(s)
}

func (s UserLogSlice) Less(i, j int) bool {
	return s[i].Msg < s[j].Msg
}

func (s UserLogSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type TestAudit struct {
	Auditer
	userLogs UserLogSlice
}

func NewTestAudit() *TestAudit {
	return &TestAudit{}
}

func (a *TestAudit) LogUserCmd(msg string, result uint) {
	a.userLogs = append(a.userLogs, UserLog{LOG_TYPE_USER_CMD, msg, result})
}

func (a *TestAudit) LogUserConfig(msg string, result bool) {
	var res uint
	if result {
		res = 1
	} else {
		res = 0
	}
	a.userLogs = append(a.userLogs, UserLog{LOG_TYPE_USER_CFG, msg, res})
}

func (a *TestAudit) GetUserLogs() UserLogSlice {
	ret := make(UserLogSlice, len(a.userLogs))
	copy(ret, a.userLogs)
	return ret
}

func (a *TestAudit) ClearUserLogs() {
	a.userLogs = UserLogSlice{}
}

func AssertUserLogSliceEqual(t *testing.T, exp, act UserLogSlice) {
	t.Helper()
	if !reflect.DeepEqual(exp, act) {
		t.Fatalf("Audit log mismatch: Expected: %v\nGot: %v\n", exp, act)
	}
}

func AssertUserLogSliceEqualSort(t *testing.T, exp, act UserLogSlice) {
	sort.Sort(exp)
	sort.Sort(act)
	AssertUserLogSliceEqual(t, exp, act)
}
