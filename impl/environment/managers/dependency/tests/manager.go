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
	"errors"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency"
	"github.com/componego/componego/impl/environment/managers/dependency/container"
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/testing/types"
)

func DependencyManagerTester[T testing.TRun[T]](
	t testing.TRun[T],
	factory func() (componego.DependencyInvoker, func(container.Container) error),
) {
	aStruct1 := &types.AStruct{
		Value: 123,
	}
	aStruct2 := &types.AStruct{
		Value: 321,
	}
	errCustom := errors.New("custom error")
	diManager, initializeManager := factory()
	diContainer, initializeContainer := container.New(5)
	require.NoError(t, initializeManager(diContainer))
	closeContainer, err := initializeContainer([]componego.Dependency{
		func() types.AInterface {
			return aStruct1
		},
		func() *types.AStruct {
			return aStruct2
		},
		func(aStruct *types.AStruct) *types.BStruct {
			return &types.BStruct{
				AStruct: aStruct,
			}
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, closeContainer())
	})

	t.Run("Invoke", func(t T) {
		t.Run("nil argument", func(t T) {
			_, err := diManager.Invoke(nil)
			require.ErrorIs(t, err, dependency.ErrNilArgument)
		})

		t.Run("not a function", func(t T) {
			var value *types.AStruct
			_, err := diManager.Invoke(value)
			require.ErrorIs(t, err, dependency.ErrNotFunction) // you should use Populate.
		})

		t.Run("variadic function", func(t T) {
			_, err := diManager.Invoke(func(_ ...types.AInterface) {})
			require.ErrorIs(t, err, dependency.ErrVariadicFunction)
		})

		t.Run("without arguments", func(t T) {
			value, err := diManager.Invoke(func() {})
			require.NoError(t, err)
			require.Nil(t, value)
		})

		t.Run("function with dependency", func(t T) {
			value, err := diManager.Invoke(func(aStruct *types.AStruct, _ *types.BStruct) int {
				return aStruct.Value
			})
			require.NoError(t, err)
			require.Equal(t, aStruct2.Value, value)
		})

		t.Run("return error as last value", func(t T) {
			value, err := diManager.Invoke(func(_ types.AInterface) error {
				return errCustom
			})
			require.ErrorIs(t, err, errCustom)
			require.Nil(t, value)

			value, err = diManager.Invoke(func(_ types.AInterface) error {
				return nil
			})
			require.NoError(t, err)
			require.Nil(t, value)

			value, err = diManager.Invoke(func(_ types.AInterface) (error, error) {
				return errCustom, nil
			})
			require.NoError(t, err) // because the error returned is not the last value.
			require.Same(t, errCustom, value)

			value, err = diManager.Invoke(func(_ types.AInterface) (float64, error) {
				return 123.456, nil
			})
			require.NoError(t, err)
			require.Equal(t, 123.456, value)
		})

		t.Run("getting a non-existent dependency", func(t T) {
			value, err := diManager.Invoke(func(_ *types.CStruct) {})
			require.ErrorIs(t, err, dependency.ErrDependencyManager)
			require.ErrorIs(t, err, container.ErrNotFoundType)
			require.Nil(t, value)
		})

		t.Run("returning a value from a function", func(t T) {
			value, err := diManager.Invoke(func(_ types.AInterface) (bool, string) {
				return false, ""
			})
			require.NoError(t, err)
			require.NotPanics(t, func() {
				// only the first and last return value is used.
				require.False(t, value.(bool))
			})
		})
	})

	t.Run("Populate", func(t T) {
		t.Run("value as struct", func(t T) {
			var value1 *types.AStruct
			require.NoError(t, diManager.Populate(&value1))
			require.Same(t, aStruct2, value1)
			require.Equal(t, aStruct2.Value, value1.Value)                                 // Visual comparison.
			require.ErrorIs(t, diManager.Populate(value1), dependency.ErrNotAllowedTarget) // missing &.

			var value2 types.AStruct // missing *.
			require.ErrorIs(t, diManager.Populate(&value2), dependency.ErrNotAllowedTarget)
		})

		t.Run("value as interface", func(t T) {
			var value types.AInterface
			require.NoError(t, diManager.Populate(&value))
			require.Equal(t, aStruct1.Value, value.(*types.AStruct).Value)
			require.ErrorIs(t, diManager.Populate(value), dependency.ErrNotAllowedTarget) // missing &.
		})

		t.Run("nil value", func(t T) {
			require.ErrorIs(t, diManager.Populate(nil), dependency.ErrNilArgument)
		})

		t.Run("not provided type", func(t T) {
			var value *types.CStruct
			require.ErrorIs(t, diManager.Populate(&value), dependency.ErrDependencyManager)
			require.ErrorIs(t, diManager.Populate(&value), container.ErrNotFoundType)
		})
	})

	t.Run("PopulateFields", func(t T) {
		t.Run("filling public and private keys with a special tag", func(t T) {
			value := &types.CStruct{}
			require.NoError(t, diManager.PopulateFields(value))
			require.Same(t, aStruct2, value.PublicField1)
			require.Same(t, aStruct2, value.AStruct)
			require.Same(t, aStruct1, value.GetPrivateField())
			require.Nil(t, value.PublicField2)
			require.Nil(t, value.IncorrectTag1)
			require.Nil(t, value.IncorrectTag2)
		})

		t.Run("type not provided", func(t T) {
			value := &types.DStruct{}
			require.ErrorIs(t, diManager.PopulateFields(value), dependency.ErrDependencyManager)
			require.ErrorIs(t, diManager.PopulateFields(value), container.ErrNotFoundType)
		})

		t.Run("nil argument", func(t T) {
			require.ErrorIs(t, diManager.PopulateFields(nil), dependency.ErrNilArgument)
		})

		t.Run("not a pointer", func(t T) {
			value1 := types.CStruct{}
			require.ErrorIs(t, diManager.PopulateFields(value1), dependency.ErrNotAllowedTarget)
			value2 := "string"
			require.ErrorIs(t, diManager.PopulateFields(&value2), dependency.ErrNotAllowedTarget)
		})
	})
}
