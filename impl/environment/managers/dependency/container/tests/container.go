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
	"reflect"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency/container"
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/testing/types"
)

func DependencyContainerTester[T testing.T](
	t testing.TRun[T],
	factory func() (container.Container, func([]componego.Dependency) error),
) {
	t.Run("basic constructor", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		// Get Value again
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err2)
		require.Same(t, reflectValue1.Interface(), reflectValue2.Interface())
	})

	t.Run("dependency as value", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			&types.AStruct{},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		// Get Value again
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err2)
		require.Same(t, reflectValue1.Interface(), reflectValue2.Interface())
	})

	t.Run("getting a value that is not provided", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			&types.AStruct{},
		}))
		_, err := c.GetValue(reflect.TypeOf((*types.BStruct)(nil)))
		require.ErrorIs(t, err, container.ErrNotFoundType)
	})

	t.Run("constructor that returns an error", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() (*types.AStruct, error) {
				return &types.AStruct{}, nil
			},
		}))
		reflectValue, err := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err)
		require.IsType(t, &types.AStruct{}, reflectValue.Interface())
	})

	t.Run("error is not the last return value", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() (error, *types.AStruct) {
				return nil, &types.AStruct{}
			},
		}), container.ErrInvalidProvidedType)
	})

	t.Run("returning multiple Values from the one constructor", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() (*types.AStruct, *types.BStruct, error) {
				return &types.AStruct{}, &types.BStruct{}, nil
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.BStruct)(nil)))
		require.NoError(t, err2)
		require.IsType(t, &types.BStruct{}, reflectValue2.Interface())
	})

	t.Run("returning values from the several constructor", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
			func() (*types.BStruct, error) {
				return &types.BStruct{}, nil
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.BStruct)(nil)))
		require.NoError(t, err2)
		require.IsType(t, &types.BStruct{}, reflectValue2.Interface())
	})

	t.Run("the constructor returns an error", func(t T) {
		_, initializer := factory()
		errCustom := errors.New("custom error")
		err := initializer([]componego.Dependency{
			func() (*types.AStruct, error) {
				return &types.AStruct{}, errCustom
			},
		})
		require.ErrorIs(t, err, errCustom)
		require.ErrorIs(t, err, container.ErrGettingDependency)
	})

	t.Run("constructor that expects other dependencies", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
			func(aStruct *types.AStruct) (*types.BStruct, error) {
				return &types.BStruct{
					AStruct: aStruct,
				}, nil
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.BStruct)(nil)))
		require.NoError(t, err2)
		require.IsType(t, &types.BStruct{}, reflectValue2.Interface())
		require.Same(t, reflectValue2.Interface().(*types.BStruct).AStruct, reflectValue1.Interface())
	})

	t.Run("the constructor expects a type that was not provided", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func(aStruct *types.AStruct) (*types.BStruct, error) {
				return &types.BStruct{
					AStruct: aStruct,
				}, nil
			},
		}), container.ErrUndeclaredDependency)
	})

	t.Run("cyclic dependencies", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func(_ *types.BStruct) *types.AStruct {
				return &types.AStruct{}
			},
			func(_ *types.AStruct) (*types.BStruct, error) {
				return &types.BStruct{}, nil
			},
		}), container.ErrCyclicDependencies)
	})

	t.Run("no dependencies", func(t T) {
		_, initializer := factory()
		require.NoError(t, initializer(nil))
		_, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{}))
		_, initializer = factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() {},
		}), container.ErrNoReturnDependencies)
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() error { return nil },
		}), container.ErrNoReturnDependencies)
		require.NoError(t, initializer([]componego.Dependency{
			func() (*types.AStruct, error) { return &types.AStruct{}, nil },
		}))
	})

	t.Run("not supported type for dependency constructor", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{nil}), container.ErrNilFactory)
		notSupportedConstructors := [...]componego.Dependency{
			nil,
			int32(1),
			func() int32 { return 1 },
			int64(1),
			func() int64 { return 1 },
			float32(1),
			func() float32 { return 1 },
			float64(1),
			func() float64 { return 1 },
			"string",
			func() string { return "" },
			' ',
			func() byte { return ' ' },
			make(chan int),
			func() chan int { return make(chan int) },
			types.AStruct{}, // because it should be a pointer.
			func() types.AStruct { return types.AStruct{} },
			(func() *int32 {
				// simple type like a pointer.
				value := int32(1)
				return &value
			})(),
			func() *int32 { return nil },
		}
		for _, constructor := range notSupportedConstructors {
			_, initializer = factory()
			require.ErrorIs(t, initializer([]componego.Dependency{constructor}), container.ErrInvalidProvidedType)
		}
	})

	t.Run("variadic constructor", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
			func(_ ...*types.AStruct) *types.BStruct {
				return &types.BStruct{}
			},
		}), container.ErrVariadicFactory)
	})

	t.Run("constructor that returns a function as a dependency", func(t T) {
		c, initializer := factory()
		someFunc := func() int { return 123 }
		require.NoError(t, initializer([]componego.Dependency{
			func() (func() int, error) {
				return someFunc, nil
			},
		}))
		reflectValue, err := c.GetValue(reflect.TypeOf(someFunc))
		require.NoError(t, err)
		require.Equal(t, reflectValue.Interface().(func() int)(), 123)
	})

	t.Run("constructor that returns an interface", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{
					Value: 123,
				}
			},
			func() types.AInterface {
				return &types.AStruct{
					Value: 321,
				}
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.IsType(t, &types.AStruct{}, reflectValue1.Interface())
		reflectValue2, err2 := c.GetValue(reflect.TypeOf((*types.AInterface)(nil)).Elem())
		require.NoError(t, err2)
		require.Equal(t, 123, reflectValue1.Interface().(*types.AStruct).Value)
		require.Equal(t, 321, reflectValue2.Interface().(*types.AStruct).Value)
	})

	t.Run("several identical return types", func(t T) {
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() (*types.AStruct, *types.AStruct) {
				return &types.AStruct{}, &types.AStruct{}
			},
		}), container.ErrSameDependencyType)
		_, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			// The return types are different in this case.
			func() (types.AInterface, *types.AStruct) {
				return &types.AStruct{}, &types.AStruct{}
			},
		}))
	})

	t.Run("rewriting dependencies", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{
					Value: 123,
				}
			},
			func() *types.AStruct {
				return &types.AStruct{
					Value: 321,
				}
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 321, reflectValue1.Interface().(*types.AStruct).Value)
		// Reverse order.
		c, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{
					Value: 321,
				}
			},
			func() *types.AStruct {
				return &types.AStruct{
					Value: 123,
				}
			},
		}))
		reflectValue1, err1 = c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 123, reflectValue1.Interface().(*types.AStruct).Value)
	})

	t.Run("rewriting a type that is added without using a constructor", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			&types.AStruct{
				Value: 123,
			},
			func() *types.AStruct {
				return &types.AStruct{
					Value: 321,
				}
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 321, reflectValue1.Interface().(*types.AStruct).Value)
		// Reverse order.
		c, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{
					Value: 321,
				}
			},
			&types.AStruct{
				Value: 123,
			},
		}))
		reflectValue1, err1 = c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 123, reflectValue1.Interface().(*types.AStruct).Value)
	})

	t.Run("rewriting dependencies with error as last return Value", func(t T) {
		c, initializer := factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{
					Value: 123,
				}
			},
			func() (*types.AStruct, error) {
				return &types.AStruct{
					Value: 321,
				}, nil
			},
		}))
		reflectValue1, err1 := c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 321, reflectValue1.Interface().(*types.AStruct).Value)
		// Reverse order.
		c, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() (*types.AStruct, error) {
				return &types.AStruct{
					Value: 321,
				}, nil
			},
			func() *types.AStruct {
				return &types.AStruct{
					Value: 123,
				}
			},
		}))
		reflectValue1, err1 = c.GetValue(reflect.TypeOf((*types.AStruct)(nil)))
		require.NoError(t, err1)
		require.Equal(t, 123, reflectValue1.Interface().(*types.AStruct).Value)
	})

	t.Run("incorrect dependency rewrite", func(t T) {
		// Return types must match the previous constructor.
		_, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() (*types.AStruct, types.AInterface) {
				return &types.AStruct{}, &types.AStruct{}
			},
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
		}), container.ErrIncorrectRewrite)
		// The number of Values returned has been increased, so there is no error.
		_, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
		}))
		// The number of return Values has been reduced, causing an error.
		_, initializer = factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
			func() *types.AStruct {
				return &types.AStruct{}
			},
		}), container.ErrIncorrectRewrite)
		// The same is true in this case.
		// The number of return Values has been reduced, causing an error.
		_, initializer = factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
			&types.AStruct{},
		}), container.ErrIncorrectRewrite)

		_, initializer = factory()
		require.ErrorIs(t, initializer([]componego.Dependency{
			func() *types.AStruct {
				return &types.AStruct{}
			},
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
			func() *types.AStruct { // wrong rewrite here.
				return &types.AStruct{}
			},
		}), container.ErrIncorrectRewrite)

		_, initializer = factory()
		require.NoError(t, initializer([]componego.Dependency{
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
			func() *types.AStruct { // this is ignored since the correct rewrite is given below.
				return &types.AStruct{}
			},
			func() (*types.AStruct, *types.BStruct) {
				return &types.AStruct{}, &types.BStruct{}
			},
		}))
	})
}

func GenerateTestFactories(countFactories int, countReturnTypes int) []componego.Dependency {
	result := make([]componego.Dependency, countFactories)
	for i := 0; i < countFactories; i++ {
		returnTypes := make([]reflect.Type, countReturnTypes)
		returnValues := make([]reflect.Value, countReturnTypes)
		for j := 0; j < countReturnTypes; j++ {
			returnType := reflect.StructOf([]reflect.StructField{
				{
					Name: fmt.Sprintf("Field_%d_%d", i, j),
					Type: reflect.TypeOf(float64(0)),
				},
			})
			returnTypes[j] = reflect.PointerTo(returnType)
			returnValues[j] = reflect.Zero(returnTypes[j])
		}
		funcType := reflect.FuncOf([]reflect.Type{}, returnTypes, false)
		result[i] = reflect.MakeFunc(funcType, func(_ []reflect.Value) []reflect.Value {
			return returnValues
		}).Interface()
	}
	return result
}
