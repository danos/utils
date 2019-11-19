// Copyright (c) 2019, AT&T Intellectual Property Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package guard_test

import (
	"fmt"
	"testing"

	"github.com/danos/utils/guard"
)

const errmsg = `Intentional Panic`

var testerr = fmt.Errorf("Test Error Message")

func letsPanicError(cause string, haveapanic bool, err error) error {
	if haveapanic {
		panic(cause)
	}
	return err
}

func letsPanicBoolError(cause string, result, haveapanic bool, err error) (bool, error) {
	if haveapanic {
		panic(cause)
	}
	return result, err
}

func letsPanicIntError(cause string, result int, haveapanic bool, err error) (int, error) {
	if haveapanic {
		panic(cause)
	}
	return result, err
}

func letsPanicStringError(cause string, result string, haveapanic bool, err error) (string, error) {
	if haveapanic {
		panic(cause)
	}
	return result, err
}

func isErrAsExpected(t *testing.T, expect, got error) {
	switch {
	case expect == nil && got != nil:
		t.Fatalf("Unexpected error message: %s", got.Error())
	case expect != nil && got == nil:
		t.Fatalf("Unexpectedly failed to return error")
	case expect != nil && got != nil:
		if expect.Error() != got.Error() {
			t.Fatalf("Expected error: %s\n but got error: %s\n",
				expect.Error(), got.Error())
		}
	}
	return

}

func isBoolAsExpected(t *testing.T, expect, got bool) {
	if expect != got {
		t.Fatalf("Returned bool not as expected:\n Expected: %t\n Got: %t\n", expect, got)
	}
}

func isIntAsExpected(t *testing.T, expect, got int) {
	if expect != got {
		t.Fatalf("Returned bool not as expected:\n Expected: %d\n Got: %d\n", expect, got)
	}
}

func isStringAsExpected(t *testing.T, expect, got string) {
	if expect != got {
		t.Fatalf("Returned bool not as expected:\n Expected: %s\n Got: %s\n", expect, got)
	}
}

// TestNoPanicErrorOnly tests that CatchPanicErrorOnly correctly returns with
// the correct error, originating from the guarded function, not a panic.
func TestNoPanicErrorOnly(t *testing.T) {
	// Check no error as expected
	err := guard.CatchPanicErrorOnly(func() error { return letsPanicError(errmsg, false, nil) })
	isErrAsExpected(t, nil, err)

	// Check expected error seen
	err = guard.CatchPanicErrorOnly(func() error { return letsPanicError(errmsg, false, testerr) })
	isErrAsExpected(t, testerr, err)

}

// TestPanicErrorOnly tests that a panic is caught, and the returned error
// originates from the panic, not from the function being guarded
func TestPanicErrorOnly(t *testing.T) {
	err := guard.CatchPanicErrorOnly(func() error { return letsPanicError(errmsg, true, nil) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)

	err = guard.CatchPanicErrorOnly(func() error { return letsPanicError(errmsg, true, testerr) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
}

// TestNoPanicBoolError tests that when no panic occurs during call of guarded
// function, the expected result and/or error is observed
func TestNoPanicBoolError(t *testing.T) {
	result, err := guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, false, false, nil) })
	isErrAsExpected(t, nil, err)
	isBoolAsExpected(t, false, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, true, false, nil) })
	isErrAsExpected(t, nil, err)
	isBoolAsExpected(t, true, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, false, false, testerr) })
	isErrAsExpected(t, testerr, err)
	isBoolAsExpected(t, false, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, true, false, testerr) })
	isErrAsExpected(t, testerr, err)
	isBoolAsExpected(t, true, result)

}

// TestNoPanicBoolError tests that when no panic occurs during call of guarded
// function, the expected result and/or error is observed
func TestPanicBoolError(t *testing.T) {
	result, err := guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, false, true, nil) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, true, true, nil) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, false, true, testerr) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)

	result, err = guard.CatchPanicBoolError(func() (bool, error) { return letsPanicBoolError(errmsg, true, true, testerr) })
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)
}

// TestCatchPanicNoPanic tests that CatchPanic correctly returns arguments of different types
// with the expected values
func TestCatchPanicNoPanic(t *testing.T) {
	b, err := guard.CatchPanic(func() (interface{}, error) { return letsPanicBoolError(errmsg, true, false, nil) })
	result, _ := b.(bool)
	isErrAsExpected(t, nil, err)
	isBoolAsExpected(t, true, result)

	b, err = guard.CatchPanic(func() (interface{}, error) { return letsPanicStringError(errmsg, "TestString", false, testerr) })
	strResult, _ := b.(string)
	isErrAsExpected(t, testerr, err)
	isStringAsExpected(t, "TestString", strResult)

	b, err = guard.CatchPanic(func() (interface{}, error) { return letsPanicIntError(errmsg, 42, false, testerr) })
	intResult, _ := b.(int)
	isErrAsExpected(t, testerr, err)
	isIntAsExpected(t, 42, intResult)
}

// TestCatchPanicWithPanic tests that CatchPanic correctly returns an error originating from a panic.
// any return values are an appropriate nil value for the type
func TestCatchPanicWithPanic(t *testing.T) {
	b, err := guard.CatchPanic(func() (interface{}, error) { return letsPanicBoolError(errmsg, true, true, nil) })
	result, _ := b.(bool)
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)

	b, err = guard.CatchPanic(func() (interface{}, error) { return letsPanicBoolError(errmsg, true, true, testerr) })
	result, _ = b.(bool)
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isBoolAsExpected(t, false, result)

	b, err = guard.CatchPanic(func() (interface{}, error) { return letsPanicIntError(errmsg, 42, true, testerr) })
	intResult, _ := b.(int)
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isIntAsExpected(t, 0, intResult)

	b, err = guard.CatchPanic(func() (interface{}, error) { return letsPanicStringError(errmsg, "Test String", true, testerr) })
	strResult, _ := b.(string)
	isErrAsExpected(t, fmt.Errorf(errmsg), err)
	isStringAsExpected(t, "", strResult)

}
