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
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing/require"
	xerrorsTests "github.com/componego/componego/libs/xerrors/tests"
)

func TestComponentManager(t *testing.T) {
	t.Run("common compatibility tests", func(t *testing.T) {
		ComponentManagerTester[*testing.T](t, component.NewManager)
	})

	t.Run("compare cycle detection error", func(t *testing.T) {
		t.Run("a component that depends on itself", func(t *testing.T) {
			compareCycleError(t, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:     "component1",
							version:  "0.0.1",
							children: []treeItem{
								// ...
							},
						},
					},
				},
			}, []string{
				"component1@0.0.1",
				"component1@0.0.1",
			})
		})

		t.Run("a component that depends on another component that depends on first one", func(t *testing.T) {
			compareCycleError(t, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
							children: []treeItem{
								{
									name:     "component1",
									version:  "0.0.1",
									children: []treeItem{
										// ...
									},
								},
							},
						},
					},
				},
			}, []string{
				"component1@0.0.1",
				"component2@0.0.1",
				"component1@0.0.1",
			})
		})

		t.Run("a component that depends on another component that depends on first one with another version", func(t *testing.T) {
			compareCycleError(t, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
							children: []treeItem{
								{
									name:    "component3",
									version: "0.0.1",
									children: []treeItem{
										{
											name:     "component2",
											version:  "0.0.2", // new version
											children: []treeItem{
												// ...
											},
										},
									},
								},
							},
						},
					},
				},
			}, []string{
				"component2@0.0.1",
				"component3@0.0.1",
				"component2@0.0.2",
			})
		})
	})
}

func TestExtractComponents(t *testing.T) {
	t.Run("application without components", func(t *testing.T) {
		app := &testApplication{}
		components, err := component.ExtractComponents(app)
		require.NoError(t, err)
		require.Len(t, components, 0)
		// Additional check to ensure that the application does not implement the componego.ApplicationComponents interface.
		if _, ok := any(app).(componego.ApplicationComponents); ok {
			t.FailNow()
		}
	})

	t.Run("application with components", func(t *testing.T) {
		appFactory := application.NewFactory("test")

		t.Run("that does not return a error", func(t *testing.T) {
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				return createTree([]treeItem{
					{
						name:    "component1",
						version: "0.0.1",
					},
				}), nil
			})
			components, err := component.ExtractComponents(appFactory.Build())
			require.NoError(t, err)
			require.Equal(t, []string{
				"component1@0.0.1",
			}, getComponentsIdentifiers(components))
		})

		t.Run("that returns a error", func(t *testing.T) {
			expectedErr := errors.New("error 1")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				return createTree([]treeItem{
					{
						name:    "component1",
						version: "0.0.1",
					},
				}), expectedErr
			})
			components, actualErr := component.ExtractComponents(appFactory.Build())
			require.ErrorIs(t, actualErr, expectedErr)
			require.ErrorIs(t, actualErr, component.ErrComponentManager)
			require.Len(t, components, 0)
		})
	})
}

type testApplication struct{}

func (t testApplication) ApplicationName() string {
	panic("not implemented")
}

func (t testApplication) ApplicationAction(_ componego.Environment, _ any) (int, error) {
	panic("not implemented")
}

func compareCycleError(t *testing.T, components []treeItem, expectedIdentifier []string) {
	manager, initializer := component.NewManager()
	err := initializer(createTree(components))
	require.ErrorIs(t, err, component.ErrCyclicDependencies)
	require.Len(t, manager.Components(), 0)
	cyclicDependencies := xerrorsTests.GetOptionsByKey(t, err, "component:cyclicDependencies")
	require.NotPanics(t, func() {
		require.Equal(t, expectedIdentifier, getComponentsIdentifiers(cyclicDependencies.([]componego.Component)))
	})
}
