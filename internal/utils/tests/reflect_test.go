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
	"reflect"
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/utils"
)

func TestIsErrorType(t *testing.T) {
	t.Run("built-in error type", func(t *testing.T) {
		errType := reflect.TypeOf((*error)(nil)).Elem()
		require.True(t, utils.IsErrorType(errType))
	})

	t.Run("custom error type", func(t *testing.T) {
		errType := reflect.TypeOf((*customError)(nil)).Elem()
		require.True(t, utils.IsErrorType(errType))
	})

	t.Run("not error type", func(t *testing.T) {
		type notAnError struct{}
		notErrType := reflect.TypeOf((*notAnError)(nil)).Elem()
		require.False(t, utils.IsErrorType(notErrType))
	})
}

func TestIndirect(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var instance any
		result := utils.Indirect(instance)
		require.Nil(t, result)
	})

	t.Run("non-pointer value", func(t *testing.T) {
		instance := 123
		result := utils.Indirect(instance)
		require.Equal(t, 123, result)
	})

	t.Run("single pointer value", func(t *testing.T) {
		value := 123
		instance := &value
		result := utils.Indirect(instance)
		require.Equal(t, 123, result)
	})

	t.Run("double pointer value", func(t *testing.T) {
		value := 123
		instance := &value
		doublePointer := &instance
		result := utils.Indirect(doublePointer)
		require.Equal(t, 123, result)
	})

	t.Run("triple pointer value", func(t *testing.T) {
		value := 123
		instance := &value
		doublePointer := &instance
		triplePointer := &doublePointer
		result := utils.Indirect(triplePointer)
		require.Equal(t, 123, result)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var instance *int
		result := utils.Indirect(instance)
		require.Nil(t, result)
	})

	t.Run("pointer to nil pointer", func(t *testing.T) {
		var value *int
		instance := &value
		result := utils.Indirect(instance)
		require.Nil(t, result)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		type testStruct struct {
			Field int
		}
		instance := &testStruct{Field: 123}
		result := utils.Indirect(instance)
		require.Equal(t, testStruct{Field: 123}, result)
	})

	t.Run("pointer to pointer to struct", func(t *testing.T) {
		type testStruct struct {
			Field int
		}
		value := &testStruct{Field: 123}
		instance := &value
		result := utils.Indirect(instance)
		require.Equal(t, testStruct{Field: 123}, result)
	})

	t.Run("pointer to nil interface", func(t *testing.T) {
		var value any
		instance := &value
		result := utils.Indirect(instance)
		require.Nil(t, result)
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		var instance any
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("empty string", func(t *testing.T) {
		instance := ""
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-empty string", func(t *testing.T) {
		instance := "hello"
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("empty slice", func(t *testing.T) {
		var instance []int
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-empty slice", func(t *testing.T) {
		instance := []int{1, 2, 3}
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("empty map", func(t *testing.T) {
		instance := map[string]int{}
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-empty map", func(t *testing.T) {
		instance := map[string]int{"key": 123}
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("nil pointer", func(t *testing.T) {
		var instance *int
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		value := 123
		instance := &value
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("pointer to empty struct", func(t *testing.T) {
		instance := &struct{}{}
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("pointer to non empty struct", func(t *testing.T) {
		instance := &struct{ x int }{x: 123}
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("empty chan", func(t *testing.T) {
		instance := make(chan int)
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-empty chan", func(t *testing.T) {
		instance := make(chan int, 1)
		instance <- 1
		require.False(t, utils.IsEmpty(instance))
	})

	t.Run("zero struct", func(t *testing.T) {
		instance := struct{ x int }{}
		require.True(t, utils.IsEmpty(instance))
	})

	t.Run("non-zero struct", func(t *testing.T) {
		instance := struct{ x int }{x: 123}
		require.False(t, utils.IsEmpty(instance))
	})
}

type customError struct {
	message string
}

func (c customError) Error() string {
	return c.message
}

var _ error = (*customError)(nil)
