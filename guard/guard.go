// Copyright (c) 2019, AT&T Intellectual Property Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package guard

import (
	"fmt"
)

func CatchPanic(fn func() (interface{}, error)) (out interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
		}
	}()
	return fn()
}

func CatchPanicErrorOnly(fn func() error) (err error) {
	_, err = CatchPanic(func() (interface{}, error) {
		e := fn()
		return nil, e
	})
	return
}

func CatchPanicBoolError(fn func() (bool, error)) (result bool, err error) {
	r, err := CatchPanic(func() (interface{}, error) {
		return fn()
	})
	result, _ = r.(bool)
	return
}
