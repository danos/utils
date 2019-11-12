// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package audit

type Auditer interface {
	LogUserCmd(msg string, result uint)
	LogUserConfig(msg string, result bool)
}

type Audit struct {
	Auditer
}

func NewAudit() *Audit {
	return &Audit{}
}

func (a Audit) LogUserCmd(msg string, result uint) {
	LogUserCmd(msg, result)
}

func (a Audit) LogUserConfig(msg string, result bool) {
	LogUserConfig(msg, result)
}
