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
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing/require"
)

func ComponentManagerTester(
	t *testing.T,
	factory func() (componego.ComponentProvider, func([]componego.Component) error),
) {
	redLevel10 := component.NewFactory("red", "red:1.0")
	redLevel10.SetComponentComponents(func() ([]componego.Component, error) {
		greenLevel20 := component.NewFactory("green", "green:2.0")
		greenLevel20.SetComponentComponents(func() ([]componego.Component, error) {
			blueLevel30 := component.NewFactory("blue", "blue:3.0")
			purpleLevel30 := component.NewFactory("purple", "purple:3.0")
			return []componego.Component{
				blueLevel30.Build(),
				purpleLevel30.Build(),
			}, nil
		})
		orangeLevel20 := component.NewFactory("orange", "orange:2.0")
		orangeLevel20.SetComponentComponents(func() ([]componego.Component, error) {
			purpleLevel31 := component.NewFactory("purple", "purple:3.1")
			return []componego.Component{
				purpleLevel31.Build(),
			}, nil
		})
		return []componego.Component{
			greenLevel20.Build(),
			orangeLevel20.Build(),
		}, nil
	})
	purpleLevel10 := component.NewFactory("purple", "purple:1.0")
	orangeLevel10 := component.NewFactory("orange", "orange:1.0")
	orangeLevel10.SetComponentComponents(func() ([]componego.Component, error) {
		purpleLevel20 := component.NewFactory("purple", "purple:2.0")
		return []componego.Component{
			purpleLevel20.Build(),
		}, nil
	})
	// run tests
	t.Run("test 1", func(t *testing.T) {
		manager, initializer := factory()
		require.NoError(t, initializer([]componego.Component{
			redLevel10.Build(),
			purpleLevel10.Build(),
			orangeLevel10.Build(),
		}))
		// todo: improve the list. We need to add lists for different manager implementations.
		require.Equal(t, []string{
			"blue:3.0",
			"purple:2.0",
			"green:2.0",
			"orange:1.0",
			"red:1.0",
		}, getUniqueIdentifiers(manager.Components()))
	})
	t.Run("test 2", func(t *testing.T) {
		manager, initializer := factory()
		require.NoError(t, initializer([]componego.Component{
			purpleLevel10.Build(),
			orangeLevel10.Build(),
			redLevel10.Build(),
		}))
		// todo: improve the list. We need to add lists for different manager implementations.
		require.Equal(t, []string{
			"blue:3.0",
			"purple:3.1",
			"green:2.0",
			"orange:2.0",
			"red:1.0",
		}, getUniqueIdentifiers(manager.Components()))
	})
	t.Run("test 3", func(t *testing.T) {
		blueLevel1 := component.NewFactory("blue", "blue:1.0")
		manager, initializer := factory()
		require.NoError(t, initializer([]componego.Component{
			redLevel10.Build(),
			blueLevel1.Build(),
		}))
		// todo: improve the list. We need to add lists for different manager implementations.
		require.Equal(t, []string{
			"purple:3.1",
			"blue:1.0",
			"green:2.0",
			"orange:2.0",
			"red:1.0",
		}, getUniqueIdentifiers(manager.Components()))
	})
	t.Run("cyclic dependencies 1", func(t *testing.T) {
		blueLevel10 := component.NewFactory("blue", "blue:1.0")
		blueLevel10.SetComponentComponents(func() ([]componego.Component, error) {
			return []componego.Component{
				redLevel10.Build(),
			}, nil
		})
		manager, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Component{
			redLevel10.Build(),
			blueLevel10.Build(),
		}), component.ErrCyclicDependencies)
		require.Len(t, manager.Components(), 0)
	})
	t.Run("cyclic dependencies 2", func(t *testing.T) {
		blueLevel10 := component.NewFactory("blue", "blue:1.0")
		blueLevel10.SetComponentComponents(func() ([]componego.Component, error) {
			return []componego.Component{
				blueLevel10.Build(),
			}, nil
		})
		manager, initializer := factory()
		require.ErrorIs(t, initializer([]componego.Component{
			redLevel10.Build(),
			blueLevel10.Build(),
		}), component.ErrCyclicDependencies)
		require.Len(t, manager.Components(), 0)
	})
	t.Run("get components", func(t *testing.T) {
		manager, initializer := factory()
		require.NoError(t, initializer([]componego.Component{
			redLevel10.Build(),
		}))
		require.Equal(t, manager.Components(), manager.Components())
		require.NotSame(t, manager.Components(), manager.Components())
	})
}

func getUniqueIdentifiers(components []componego.Component) []string {
	result := make([]string, len(components))
	for i, componentItem := range components {
		result[i] = componentItem.ComponentVersion()
	}
	return result
}
