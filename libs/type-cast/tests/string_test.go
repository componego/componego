/*
Copyright 2024-present Volodymyr Konstanchuk and contributors

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
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/type-cast"
)

func TestToString(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		result, err := type_cast.ToString(nil)
		require.NoError(t, err)
		require.Equal(t, "", result)
	})

	t.Run("bool value", func(t *testing.T) {
		result, err := type_cast.ToString(true)
		require.NoError(t, err)
		require.Equal(t, "true", result)
	})

	t.Run("int value", func(t *testing.T) {
		result, err := type_cast.ToString(123)
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("int8 value", func(t *testing.T) {
		result, err := type_cast.ToString(int8(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("int16 value", func(t *testing.T) {
		result, err := type_cast.ToString(int16(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("int32 value", func(t *testing.T) {
		result, err := type_cast.ToString(int32(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("int64 value", func(t *testing.T) {
		result, err := type_cast.ToString(int64(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("uint value", func(t *testing.T) {
		result, err := type_cast.ToString(uint(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("uint8 value", func(t *testing.T) {
		result, err := type_cast.ToString(uint8(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("uint16 value", func(t *testing.T) {
		result, err := type_cast.ToString(uint16(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("uint32 value", func(t *testing.T) {
		result, err := type_cast.ToString(uint32(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("uint64 value", func(t *testing.T) {
		result, err := type_cast.ToString(uint64(123))
		require.NoError(t, err)
		require.Equal(t, "123", result)
	})

	t.Run("float32 value", func(t *testing.T) {
		result, err := type_cast.ToString(float32(123.123))
		require.NoError(t, err)
		require.Equal(t, "123.123", result)
	})

	t.Run("float64 value", func(t *testing.T) {
		// noinspection ALL
		result, err := type_cast.ToString(float64(123.123))
		require.NoError(t, err)
		require.Equal(t, "123.123", result)
	})

	t.Run("string value", func(t *testing.T) {
		result, err := type_cast.ToString("hello")
		require.NoError(t, err)
		require.Equal(t, "hello", result)
	})

	t.Run("unsupported value type", func(t *testing.T) {
		_, err := type_cast.ToString([]int{1, 2, 3})
		require.Error(t, err)
	})
}
