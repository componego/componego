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

func TestToBool(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		result, err := type_cast.ToBool(nil)
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("bool value true", func(t *testing.T) {
		result, err := type_cast.ToBool(true)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("bool value false", func(t *testing.T) {
		result, err := type_cast.ToBool(false)
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("int non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(123)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("int zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(0)
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("int8 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int8(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("int8 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int8(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("int16 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int16(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("int16 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int16(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("int32 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int32(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("int32 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int32(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("int64 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int64(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("int64 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(int64(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("uint non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("uint zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("uint8 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint8(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("uint8 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint8(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("uint16 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint16(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("uint16 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint16(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("uint32 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint32(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("uint32 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint32(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("uint64 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint64(123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("uint64 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(uint64(0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("float32 non-zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(float32(123.123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("float32 zero value", func(t *testing.T) {
		result, err := type_cast.ToBool(float32(0.0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("float64 non-zero value", func(t *testing.T) {
		// noinspection ALL
		result, err := type_cast.ToBool(float64(123.123))
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("float64 zero value", func(t *testing.T) {
		// noinspection ALL
		result, err := type_cast.ToBool(float64(0.0))
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("string true value", func(t *testing.T) {
		result, err := type_cast.ToBool("true")
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("string false value", func(t *testing.T) {
		result, err := type_cast.ToBool("false")
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("string 1 value", func(t *testing.T) {
		result, err := type_cast.ToBool("1")
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("string 0 value", func(t *testing.T) {
		result, err := type_cast.ToBool("0")
		require.NoError(t, err)
		require.False(t, result)
	})

	t.Run("string invalid value", func(t *testing.T) {
		_, err := type_cast.ToBool("invalid")
		require.Error(t, err)
	})

	t.Run("struct value invalid", func(t *testing.T) {
		_, err := type_cast.ToBool(struct{}{})
		require.Error(t, err)
	})
}
