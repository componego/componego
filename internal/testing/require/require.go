/*
Copyright 2024 Volodymyr Konstanchuk and the Componego Framework contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package require

import (
	"testing"
)

// TODO: it is necessary to implement this using code generation.

// Equal is a proxy function for require.Equal.
func Equal(t *testing.T, expected any, actual any, msgAndArgs ...any) {
	t.Helper()
	call("Equal", t, expected, actual, msgAndArgs)
}

// Same is a proxy function for require.Same.
func Same(t *testing.T, expected any, actual any, msgAndArgs ...any) {
	t.Helper()
	call("Same", t, expected, actual, msgAndArgs)
}

// NotSame is a proxy function for require.NotSame.
func NotSame(t *testing.T, expected any, actual any, msgAndArgs ...any) {
	t.Helper()
	call("NotSame", t, expected, actual, msgAndArgs)
}

// Implements is a proxy function for require.Implements.
func Implements(t *testing.T, interfaceObject any, object any, msgAndArgs ...any) {
	t.Helper()
	call("Implements", t, interfaceObject, object, msgAndArgs)
}

// IsType is a proxy function for require.IsType.
func IsType(t *testing.T, expectedType any, object any, msgAndArgs ...any) {
	t.Helper()
	call("IsType", t, expectedType, object, msgAndArgs)
}

// ErrorIs is a proxy function for require.ErrorIs.
func ErrorIs(t *testing.T, err error, target error, msgAndArgs ...any) {
	t.Helper()
	call("ErrorIs", t, err, target, msgAndArgs)
}

// EqualError is a proxy function for require.EqualError.
func EqualError(t *testing.T, theError error, errString string, msgAndArgs ...any) {
	t.Helper()
	call("EqualError", t, theError, errString, msgAndArgs)
}

// Error is a proxy function for require.Error.
func Error(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	call("Error", t, err, msgAndArgs)
}

// NoError is a proxy function for require.NoError.
func NoError(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	call("NoError", t, err, msgAndArgs)
}

// True is a proxy function for require.True.
func True(t *testing.T, value bool, msgAndArgs ...any) {
	t.Helper()
	call("True", t, value, msgAndArgs)
}

// False is a proxy function for require.False.
func False(t *testing.T, value bool, msgAndArgs ...any) {
	t.Helper()
	call("False", t, value, msgAndArgs)
}

// Len is a proxy function for require.Len
func Len(t *testing.T, object any, length int, msgAndArgs ...any) {
	t.Helper()
	call("Len", t, object, length, msgAndArgs)
}

// Nil is a proxy function for require.Nil.
func Nil(t *testing.T, object any, msgAndArgs ...any) {
	t.Helper()
	call("Nil", t, object, msgAndArgs)
}

// NotNil is a proxy function for require.NotNil.
func NotNil(t *testing.T, object any, msgAndArgs ...any) {
	t.Helper()
	call("NotNil", t, object, msgAndArgs)
}

// Panics is a proxy function for require.Panics.
func Panics(t *testing.T, fn func(), msgAndArgs ...any) {
	t.Helper()
	call("Panics", t, fn, msgAndArgs)
}

// NotPanics is a proxy function for require.NotPanics.
func NotPanics(t *testing.T, fn func(), msgAndArgs ...any) {
	t.Helper()
	call("NotPanics", t, fn, msgAndArgs)
}

// FailNow is a proxy function for require.FailNow.
func FailNow(t *testing.T, failureMessage string, msgAndArgs ...any) {
	t.Helper()
	call("FailNow", t, failureMessage, msgAndArgs)
}

// Regexp is a proxy function for require.Regexp.
func Regexp(t *testing.T, regexp any, str any, msgAndArgs ...any) {
	t.Helper()
	call("Regexp", t, regexp, str, msgAndArgs)
}
