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

package tests

import (
	"errors"
	"testing"

	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/testing/types"
)

func TestDriver(t *testing.T) {
	DriverTester[*testing.T](t, func() driver.Driver {
		return driver.New(nil)
	})
}

func TestErrorRecoveryOnStop(t *testing.T) {
	t.Run("nil recovery without the previous error", func(t *testing.T) {
		resultErr := driver.ErrorRecoveryOnStop(nil, nil)
		require.NoError(t, resultErr)
	})

	t.Run("nil recovery with the previous error", func(t *testing.T) {
		prevErr := errors.New("previous error")
		resultErr := driver.ErrorRecoveryOnStop(nil, prevErr)
		require.NotErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorIs(t, resultErr, prevErr)
	})

	t.Run("recovery as an error without the previous error", func(t *testing.T) {
		panicErr := errors.New("panic occurred")
		resultErr := driver.ErrorRecoveryOnStop(panicErr, nil)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorIs(t, resultErr, panicErr)
	})

	t.Run("recovery as a string without the previous error", func(t *testing.T) {
		resultErr := driver.ErrorRecoveryOnStop("panic occurred", nil)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorContains(t, resultErr, "panic occurred")
	})

	t.Run("recovery as a custom string without the previous error", func(t *testing.T) {
		resultErr := driver.ErrorRecoveryOnStop(types.CustomString("panic occurred"), nil)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorContains(t, resultErr, "panic occurred")
	})

	t.Run("recovery as a custom type without the previous error", func(t *testing.T) {
		resultErr := driver.ErrorRecoveryOnStop(types.AStruct{}, nil)
		require.ErrorIs(t, resultErr, driver.ErrUnknownPanic)
	})

	t.Run("recovery as an error with the previous error", func(t *testing.T) {
		prevErr := errors.New("previous error")
		panicErr := errors.New("panic occurred")
		resultErr := driver.ErrorRecoveryOnStop(panicErr, prevErr)
		require.ErrorIs(t, resultErr, prevErr)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
	})

	t.Run("recovery as a string with the previous error", func(t *testing.T) {
		prevErr := errors.New("previous error")
		resultErr := driver.ErrorRecoveryOnStop("panic occurred", prevErr)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorIs(t, resultErr, prevErr)
		require.ErrorContains(t, resultErr, "panic occurred")
	})

	t.Run("recovery as a custom string with the previous error", func(t *testing.T) {
		prevErr := errors.New("previous error")
		resultErr := driver.ErrorRecoveryOnStop(types.CustomString("panic occurred"), prevErr)
		require.ErrorIs(t, resultErr, driver.ErrPanic)
		require.ErrorIs(t, resultErr, prevErr)
		require.ErrorContains(t, resultErr, "panic occurred")
	})

	t.Run("recovery as a custom type with the previous error", func(t *testing.T) {
		prevErr := errors.New("previous error")
		resultErr := driver.ErrorRecoveryOnStop(types.AStruct{}, prevErr)
		require.ErrorIs(t, resultErr, driver.ErrUnknownPanic)
		require.ErrorIs(t, resultErr, prevErr)
	})
}
