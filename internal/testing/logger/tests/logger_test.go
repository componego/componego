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
	"bytes"
	"fmt"
	"testing"

	"github.com/componego/componego/internal/testing/logger"
	"github.com/componego/componego/internal/testing/require"
)

func TestLogData(t *testing.T) {
	t.Run("nil writer", func(t *testing.T) {
		require.NotPanics(t, func() {
			logger.LogData(nil, "hello")
		})
	})

	t.Run("single message", func(t *testing.T) {
		var buffer bytes.Buffer
		logger.LogData(&buffer, "hello")
		expected := fmt.Sprintln("hello")
		require.Equal(t, expected, buffer.String())
	})

	t.Run("multiple messages", func(t *testing.T) {
		var buffer bytes.Buffer
		logger.LogData(&buffer, "hello", 123, true)
		expected := fmt.Sprintln("hello", 123, true)
		require.Equal(t, expected, buffer.String())
	})

	t.Run("empty messages", func(t *testing.T) {
		var buffer bytes.Buffer
		logger.LogData(&buffer)
		expected := fmt.Sprintln()
		require.Equal(t, expected, buffer.String())
	})

	t.Run("error on write", func(t *testing.T) {
		writer := &errorWriter{}
		require.Panics(t, func() {
			logger.LogData(writer, "hello")
		})
	})
}

func TestExpectedLogData(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		result := logger.ExpectedLogData(nil)
		expected := fmt.Sprintln(nil)
		require.Equal(t, expected, result)
	})

	t.Run("non-slice value", func(t *testing.T) {
		result := logger.ExpectedLogData("hello")
		expected := fmt.Sprintln("hello")
		require.Equal(t, expected, result)
	})

	t.Run("slice of strings", func(t *testing.T) {
		result := logger.ExpectedLogData([]string{"a", "b", "c"})
		expected := fmt.Sprintln([]any{"a", "b", "c"}...)
		require.Equal(t, expected, result)
	})

	t.Run("array of integers", func(t *testing.T) {
		result := logger.ExpectedLogData([...]int{1, 2, 3})
		expected := fmt.Sprintln([]any{1, 2, 3}...)
		require.Equal(t, expected, result)
	})

	t.Run("mixed values", func(t *testing.T) {
		result := logger.ExpectedLogData(nil, "hello", []int{1, 2, 3}, [...]string{"x", "y"})
		expected := fmt.Sprintln(nil) + fmt.Sprintln("hello") +
			fmt.Sprintln([]any{1, 2, 3}...) + fmt.Sprintln([]any{"x", "y"}...)
		require.Equal(t, expected, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		result := logger.ExpectedLogData([]int{})
		expected := fmt.Sprintln([]any{}...)
		require.Equal(t, expected, result)
	})

	t.Run("empty array", func(t *testing.T) {
		result := logger.ExpectedLogData([...]int{})
		expected := fmt.Sprintln([]any{}...)
		require.Equal(t, expected, result)
	})
}

// errorWriter is an io.Writer implementation that always returns an error.
type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("error")
}
