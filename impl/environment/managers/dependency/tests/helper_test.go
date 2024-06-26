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
	"fmt"
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/environment/managers/dependency"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/testing/types"
	"github.com/componego/componego/tests/runner"
)

func TestGetDependencyWithAndWithoutPanic(t *testing.T) {
	origValue := &types.AStruct{
		Value: 123,
	}
	appFactory := application.NewFactory("Test Application")
	appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
		return []componego.Dependency{
			origValue,
			func() types.AInterface {
				return origValue
			},
		}, nil
	})
	env, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), nil)
	t.Cleanup(cancelEnv)

	t.Run("get valid dependency type", func(t *testing.T) {
		value, err := dependency.Get[*types.AStruct](env)
		require.NoError(t, err)
		require.Same(t, origValue, value)
		require.NotPanics(t, func() {
			value = dependency.GetOrPanic[*types.AStruct](env)
			require.Same(t, origValue, value)
		})
	})

	t.Run("get invalid dependency type", func(t *testing.T) {
		_, err := dependency.Get[*types.BStruct](env)
		require.ErrorIs(t, err, dependency.ErrDependencyManager)
		require.Panics(t, func() {
			dependency.GetOrPanic[*types.BStruct](env)
		})
	})

	t.Run("get valid dependency as interface", func(t *testing.T) {
		value, err := dependency.Get[types.AInterface](env)
		require.NoError(t, err)
		require.Same(t, origValue, value)
		// Visual comparison.
		require.Equal(t, origValue.Value, value.(*types.AStruct).Value)
	})
}

func TestInvokeFunctionWithAndWithoutPanic(t *testing.T) {
	origValue := &types.AStruct{
		Value: 123,
	}
	errCustom := errors.New("custom error")
	appFactory := application.NewFactory("Test Application")
	appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
		return []componego.Dependency{
			origValue,
		}, nil
	})
	env, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), nil)
	t.Cleanup(cancelEnv)

	t.Run("invoke the function with dependency", func(t *testing.T) {
		value, err := dependency.Invoke[int](func(aStruct *types.AStruct) int {
			return aStruct.Value
		}, env)
		require.NoError(t, err)
		require.Equal(t, origValue.Value, value)
		require.NotPanics(t, func() {
			value = dependency.InvokeOrPanic[int](func(aStruct *types.AStruct) int {
				return aStruct.Value
			}, env)
		})
		require.Equal(t, origValue.Value, value)
	})

	t.Run("invoke the function without dependencies", func(t *testing.T) {
		value, err := dependency.Invoke[int](func() int {
			return 321
		}, env)
		require.NoError(t, err)
		require.Equal(t, 321, value)
		require.NotPanics(t, func() {
			value = dependency.InvokeOrPanic[int](func() int {
				return 321
			}, env)
			require.Equal(t, 321, value)
		})
	})

	t.Run("invoke with the incorrect type", func(t *testing.T) {
		_, err := dependency.Invoke[int]("not a function", env)
		require.ErrorIs(t, err, dependency.ErrNotFunction)
		require.Panics(t, func() {
			dependency.InvokeOrPanic[int]("not a function", env)
		})
	})

	t.Run("invoke the function with the wrong return type", func(t *testing.T) {
		_, err := dependency.Invoke[float32](func() int {
			return 321
		}, env)
		require.EqualError(t, err, fmt.Sprintf("could not convert the returned value to type %T", float32(1)))
		require.Panics(t, func() {
			dependency.InvokeOrPanic[float32](func() int {
				return 321
			}, env)
		})
	})

	t.Run("invoke the function with a missing dependency type", func(t *testing.T) {
		_, err := dependency.Invoke[float32](func(_ *types.BStruct) int {
			return 321
		}, env)
		require.ErrorIs(t, err, dependency.ErrDependencyManager)
		require.Panics(t, func() {
			dependency.InvokeOrPanic[float32](func(_ *types.BStruct) int {
				return 321
			}, env)
		})
	})

	t.Run("invoke the variadic function with the valid dependency type", func(t *testing.T) {
		_, err := dependency.Invoke[float32](func(_ ...*types.AStruct) int {
			return 321
		}, env)
		require.ErrorIs(t, err, dependency.ErrVariadicFunction)
		require.Panics(t, func() {
			dependency.InvokeOrPanic[float32](func(_ ...*types.AStruct) int {
				return 321
			}, env)
		})
	})

	t.Run("invoke the function that returns an error as the last value", func(t *testing.T) {
		_, err := dependency.Invoke[int64](func(_ *types.AStruct) error {
			return errCustom
		}, env)
		require.ErrorIs(t, err, errCustom)
		require.ErrorIs(t, err, dependency.ErrDependencyManager)
		require.Panics(t, func() {
			dependency.InvokeOrPanic[int64](func(_ *types.AStruct) error {
				return errCustom
			}, env)
		})

		_, err = dependency.Invoke[int64](func(_ *types.AStruct) (float64, error) {
			return 123, errCustom
		}, env)
		require.ErrorIs(t, err, errCustom)
		require.ErrorIs(t, err, dependency.ErrDependencyManager)
		require.Panics(t, func() {
			dependency.InvokeOrPanic[int64](func(_ *types.AStruct) (float64, error) {
				return 123, errCustom
			}, env)
		})

		_, err = dependency.Invoke[int64](func(_ *types.AStruct) (float64, error) {
			return 123, nil
		}, env)
		errString := fmt.Sprintf("could not convert the returned value to type %T", int64(1))
		require.EqualError(t, err, errString)
		require.PanicsWithError(t, errString, func() {
			dependency.InvokeOrPanic[int64](func(_ *types.AStruct) (float64, error) {
				return 123, nil
			}, env)
		})

		_, err = dependency.Invoke[error](func(_ *types.AStruct) (error, error) {
			return errCustom, nil
		}, env)
		require.NoError(t, err) // because the error returned is not the last value.
	})
}
