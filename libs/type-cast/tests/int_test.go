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
	"math"
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/type-cast"
)

func TestToInt64(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		result, err := type_cast.ToInt64(nil)
		require.NoError(t, err)
		require.Equal(t, int64(0), result)
	})

	t.Run("bool value true", func(t *testing.T) {
		result, err := type_cast.ToInt64(true)
		require.NoError(t, err)
		require.Equal(t, int64(1), result)
	})

	t.Run("bool value false", func(t *testing.T) {
		result, err := type_cast.ToInt64(false)
		require.NoError(t, err)
		require.Equal(t, int64(0), result)
	})

	t.Run("int value", func(t *testing.T) {
		result, err := type_cast.ToInt64(123)
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("int8 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(int8(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("int16 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(int16(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("int32 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(int32(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("int64 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(int64(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("uint value", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("uint8 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint8(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("uint16 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint16(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("uint32 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint32(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("uint64 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint64(123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("float32 value", func(t *testing.T) {
		result, err := type_cast.ToInt64(float32(123.123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("float64 value", func(t *testing.T) {
		// noinspection ALL
		result, err := type_cast.ToInt64(float64(123.123))
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("string value", func(t *testing.T) {
		result, err := type_cast.ToInt64("123")
		require.NoError(t, err)
		require.Equal(t, int64(123), result)
	})

	t.Run("string value invalid", func(t *testing.T) {
		_, err := type_cast.ToInt64("invalid")
		require.Error(t, err)
	})

	t.Run("struct value invalid", func(t *testing.T) {
		_, err := type_cast.ToInt64(struct{}{})
		require.Error(t, err)
	})

	t.Run("uint input above MaxInt64", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint(math.MaxInt64 + 1))
		require.Error(t, err)
		require.Equal(t, int64(0), result)
	})

	t.Run("uint64 input below MaxInt64", func(t *testing.T) {
		result, err := type_cast.ToInt64(uint64(math.MaxInt64 + 1))
		require.Error(t, err)
		require.Equal(t, int64(0), result)
	})
}
