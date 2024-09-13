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

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/require"
)

type (
	managerFactory = func() (componego.ComponentProvider, func([]componego.Component) error)
)

func ComponentManagerTester[T testing.TRun[T]](
	t T,
	factory managerFactory,
) {
	t.Run("components without dependent components", func(t T) {
		initAndCompareComponents(t, factory, []treeItem{
			{
				name:    "component1",
				version: "0.0.1",
			},
			{
				name:    "component2",
				version: "0.0.2",
			},
			{
				name:    "component3",
				version: "0.0.3",
			},
		}, []string{
			"component1@0.0.1",
			"component2@0.0.2",
			"component3@0.0.3",
		}, nil)
	})

	t.Run("empty component list", func(t T) {
		initAndCompareComponents(t, factory, nil, nil, nil)
	})

	t.Run("components with dependent components", func(t T) {
		t.Run("one level", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
						},
						{
							name:    "component5",
							version: "0.0.1",
						},
					},
				},
				{
					name:    "component6",
					version: "0.0.1",
				},
			}, []string{
				"component2@0.0.1",
				"component5@0.0.1",
				"component1@0.0.1",
				"component6@0.0.1",
			}, nil)
		})

		t.Run("two levels", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
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
								},
								{
									name:    "component4",
									version: "0.0.1",
								},
							},
						},
						{
							name:    "component5",
							version: "0.0.1",
						},
					},
				},
				{
					name:    "component6",
					version: "0.0.1",
				},
			}, []string{
				"component3@0.0.1",
				"component4@0.0.1",
				"component2@0.0.1",
				"component5@0.0.1",
				"component1@0.0.1",
				"component6@0.0.1",
			}, nil)
		})

		t.Run("three levels", func(t T) {
			// This is the same as two levels test, but with a new components.
			initAndCompareComponents(t, factory, []treeItem{
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
											name:    "component7", // new component
											version: "0.0.1",
										},
									},
								},
								{
									name:    "component4",
									version: "0.0.1",
								},
							},
						},
						{
							name:    "component5",
							version: "0.0.1",
						},
					},
				},
				{
					name:    "component6",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component9", // new component
							version: "0.0.1",
						},
					},
				},
			}, []string{
				"component7@0.0.1", // new value
				"component3@0.0.1",
				"component4@0.0.1",
				"component2@0.0.1",
				"component5@0.0.1",
				"component1@0.0.1",
				"component9@0.0.1", // new value
				"component6@0.0.1",
			}, nil)
		})
	})

	t.Run("replace components", func(t T) {
		t.Run("simple replace", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
				},
				{
					name:    "component1", // replace component1@0.0.1
					version: "0.0.2",
				},
			}, []string{
				"component1@0.0.2",
			}, nil)
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.2",
				},
				{
					name:    "component1", // replace component1@0.0.2
					version: "0.0.1",
				},
			}, []string{
				"component1@0.0.1",
			}, nil)
		})

		t.Run("replace with children", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
						},
					},
				},
				{
					name:    "component1", // replace component1@0.0.1
					version: "0.0.2",
				},
			}, []string{
				"component1@0.0.2",
			}, nil)
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.2",
				},
				{
					name:    "component1", // replace
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
						},
						{
							name:    "component2", // replace component2@0.0.1
							version: "0.0.2",
						},
					},
				},
			}, []string{
				"component2@0.0.2",
				"component1@0.0.1",
			}, nil)
		})

		t.Run("replace child components", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.2",
						},
					},
				},
				{
					name:    "component2", // replace component2@0.0.2
					version: "0.0.3",
				},
			}, []string{
				"component2@0.0.3",
				"component1@0.0.1",
			}, nil)
			initAndCompareComponents(t, factory, []treeItem{
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
								},
							},
						},
						{
							name:    "component3", // replace component3@0.0.1
							version: "0.0.2",
							children: []treeItem{
								{
									name:    "component4",
									version: "0.0.1",
								},
							},
						},
					},
				},
				{
					name:    "component4", // replace component4@0.0.1
					version: "0.0.2",
				},
				{
					name:    "component5",
					version: "0.0.1",
				},
			}, []string{
				"component4@0.0.2",
				"component3@0.0.2",
				"component2@0.0.1",
				"component1@0.0.1",
				"component5@0.0.1",
			}, nil)
		})
	})

	t.Run("cycle detection", func(t T) {
		t.Run("with same version", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component1",
							version: "0.0.1",
						},
					},
				},
			}, nil, component.ErrCyclicDependencies)
		})

		t.Run("with different version", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component1",
							version: "0.0.2",
						},
					},
				},
			}, nil, component.ErrCyclicDependencies)
		})

		t.Run("with different levels", func(t T) {
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
							children: []treeItem{
								{
									name:    "component1",
									version: "0.0.3",
								},
							},
						},
					},
				},
			}, nil, component.ErrCyclicDependencies)
		})

		t.Run("with different levels after component replace", func(t T) {
			// We replace component1@0.0.1 to component1@0.0.4 where there is no cycle.
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
							children: []treeItem{
								{
									name:    "component1",
									version: "0.0.3",
								},
							},
						},
					},
				},
				{
					name:    "component1",
					version: "0.0.4",
				},
			}, []string{
				"component1@0.0.4",
			}, nil)
			// We replace component1@0.0.4 to component1@0.0.1 where there is cycle.
			initAndCompareComponents(t, factory, []treeItem{
				{
					name:    "component1",
					version: "0.0.4",
				},
				{
					name:    "component1",
					version: "0.0.1",
					children: []treeItem{
						{
							name:    "component2",
							version: "0.0.1",
							children: []treeItem{
								{
									name:    "component1",
									version: "0.0.3",
								},
							},
						},
					},
				},
			}, nil, component.ErrCyclicDependencies)
		})
	})

	t.Run("return error instead of component list", func(t T) {
		expectedErr := errors.New("error")
		componentFactory := component.NewFactory("component1", "0.0.1")
		componentFactory.SetComponentComponents(func() ([]componego.Component, error) {
			return nil, expectedErr
		})
		manager, initializer := factory()
		actualErr := initializer([]componego.Component{
			componentFactory.Build(),
		})
		require.ErrorIs(t, actualErr, expectedErr)
		require.Len(t, manager.Components(), 0)
	})
}

type treeItem struct {
	name     string
	version  string
	children []treeItem
}

func createTree(items []treeItem) []componego.Component {
	result := make([]componego.Component, len(items))
	for index, item := range items {
		componentFactory := component.NewFactory(item.name, item.version)
		// We support compatibility with older versions of the language.
		if savedItem := item; len(savedItem.children) > 0 {
			componentFactory.SetComponentComponents(func() ([]componego.Component, error) {
				return createTree(savedItem.children), nil
			})
		}
		result[index] = componentFactory.Build()
	}
	return result
}

func getComponentsIdentifiers(components []componego.Component) []string {
	result := make([]string, len(components))
	for i, componentItem := range components {
		result[i] = componentItem.ComponentIdentifier() + "@" + componentItem.ComponentVersion()
	}
	return result
}

func initAndCompareComponents(t testing.T, factory managerFactory, components []treeItem, expectedIdentifier []string, expectedError error) {
	manager, initializer := factory()
	err := initializer(createTree(components))
	if expectedError != nil {
		require.ErrorIs(t, err, expectedError)
		require.Len(t, manager.Components(), 0)
		return
	}
	require.NoError(t, err)
	if len(expectedIdentifier) == 0 && len(manager.Components()) == 0 {
		return
	}
	require.Equal(t, expectedIdentifier, getComponentsIdentifiers(manager.Components()))
	require.Equal(t,
		getComponentsIdentifiers(manager.Components()),
		getComponentsIdentifiers(manager.Components()),
	)
}
