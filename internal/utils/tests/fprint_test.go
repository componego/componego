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
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/utils"
)

// noinspection SpellCheckingInspection
func TestFprint(t *testing.T) {
	t.Run("single string", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprint(&buffer, "hello")
		require.Equal(t, "hello", buffer.String())
	})

	t.Run("different input types", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprint(&buffer, "hello", 123, " ", '1')
		require.Equal(t, "hello123 49", buffer.String())
	})

	t.Run("empty input", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprint(&buffer)
		require.Equal(t, buffer.Len(), 0)
	})

	t.Run("error case", func(t *testing.T) {
		require.Panics(t, func() {
			utils.Fprint(&errorWriter{}, "this will fail")
		})
	})
}

// noinspection SpellCheckingInspection
func TestFprintln(t *testing.T) {
	t.Run("single string", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintln(&buffer, "hello")
		require.Equal(t, fmt.Sprintln("hello"), buffer.String())
	})

	t.Run("different input types", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintln(&buffer, "hello", 123, " ", '1')
		require.Equal(t, fmt.Sprintln("hello123 49"), buffer.String())
	})

	t.Run("empty input", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintln(&buffer)
		require.Equal(t, fmt.Sprintln(), buffer.String())
	})

	t.Run("error case", func(t *testing.T) {
		require.Panics(t, func() {
			utils.Fprintln(&errorWriter{}, "this will fail")
		})
	})
}

func TestFprintf(t *testing.T) {
	t.Run("simple formatting", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintf(&buffer, "hello %s", "world")
		require.Equal(t, "hello world", buffer.String())
	})

	t.Run("multiple format specifiers", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintf(&buffer, "integer: %d, string: %s", 123, "test")
		require.Equal(t, "integer: 123, string: test", buffer.String())
	})

	t.Run("empty format string", func(t *testing.T) {
		var buffer bytes.Buffer
		utils.Fprintf(&buffer, "")
		require.Equal(t, buffer.Len(), 0)
	})

	t.Run("error case", func(t *testing.T) {
		require.Panics(t, func() {
			utils.Fprintf(&errorWriter{}, "this will fail")
		})
	})
}

// Custom io.Writer that always returns an error.
type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("write error")
}

var _ io.Writer = (*errorWriter)(nil)
