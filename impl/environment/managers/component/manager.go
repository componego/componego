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

package component

import (
	"sort"

	"github.com/componego/componego"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/xerrors"
)

var (
	ErrComponentManager   = xerrors.New("error inside component manager", "E0410")
	ErrCyclicDependencies = ErrComponentManager.WithMessage("cycle detected in dependencies between components", "E0411")
)

type manager struct {
	components []componego.Component
}

func NewManager() (componego.ComponentProvider, func(components []componego.Component) error) {
	m := &manager{}
	return m, m.initialize
}

func (m *manager) Components() []componego.Component {
	if len(m.components) == 0 {
		return nil
	}
	// Copy the slice to avoid modification from outside the manager.
	return utils.Copy(m.components)
}

func (m *manager) initialize(components []componego.Component) error {
	if len(components) == 0 {
		return nil
	}
	componentStack := make([]*stackItem, 0, len(components)*2)
	componentMap := make(map[string]*stackItem, len(components)*2)
	for _, component := range components {
		componentStack = append(componentStack, &stackItem{
			identifier: component.ComponentIdentifier(),
			component:  component,
		})
	}
	for position := 0; len(componentStack) > 0; position++ {
		item := componentStack[len(componentStack)-1]
		componentStack = componentStack[:len(componentStack)-1]
		// You can rewrite the components using component ID.
		if oldItem, ok := componentMap[item.identifier]; ok {
			parent := item.parent
			for parent != nil {
				if parent.identifier != oldItem.identifier {
					parent = parent.parent
					continue
				}
				cyclicDependencyProvider := getCyclicDependencyProvider(item)
				return ErrCyclicDependencies.WithOptions("E0412",
					xerrors.NewCallableOption("component:cyclicDependencies", cyclicDependencyProvider),
				)
			}
			// This prevents visiting components that make no sense to visit.
			// Because this component has already been fully received.
			// Because we are checking them from the end.
			continue
		}
		componentMap[item.identifier] = item
		item.position = position
		if component, ok := item.component.(componego.ComponentComponents); ok {
			components, err := component.ComponentComponents()
			if err != nil {
				return err
			}
			children := make([]string, 0, len(components))
			for _, component := range components {
				// It is important for each new object to call this function only once during component initialization.
				identifier := component.ComponentIdentifier()
				componentStack = append(componentStack, &stackItem{
					identifier: identifier,
					component:  component,
					parent:     item,
				})
				children = append(children, identifier)
			}
			item.children = children
		}
	}
	// We add data to the previous slice, which is empty, to save memory.
	for _, item := range componentMap {
		componentStack = append(componentStack, item)
	}
	m.components = getSortedComponents(componentStack, componentMap)
	return nil
}

func ExtractComponents(app componego.Application) ([]componego.Component, error) {
	if app, ok := app.(componego.ApplicationComponents); ok {
		components, err := app.ApplicationComponents()
		if err == nil {
			return components, nil
		}
		return nil, ErrComponentManager.WithError(err, "E0413",
			xerrors.NewOption("componego:component:application", app),
		)
	}
	return nil, nil
}

func getCyclicDependencyProvider(stackItem *stackItem) func() any {
	return func() any {
		components := make([]componego.Component, 0, 5)
		components = append(components, stackItem.component)
		for parent := stackItem.parent; parent != nil; parent = parent.parent {
			components = append(components, parent.component)
			if parent.identifier == stackItem.identifier {
				break
			}
		}
		utils.Reverse(components)
		return components
	}
}

func getSortedComponents(componentStack []*stackItem, componentMap map[string]*stackItem) []componego.Component {
	// The map does not guarantee the order,
	// but it is important for us that the order in which the components are initialized is always the same.
	sort.Slice(componentStack, func(i, j int) bool {
		return componentStack[i].position < componentStack[j].position
	})
	// However, even this cannot guarantee that the order of the components will be correct.
	// We do additional sorting based on dependencies between components.
	result := make([]componego.Component, 0, len(componentMap))
	// Tracks items that have been visited during traversal to prevent revisiting.
	visited := make(map[string]bool, len(componentMap))
	// Tracks items that have been fully processed.
	processed := make(map[string]bool, len(componentMap))
	for isProcessed := true; len(componentStack) > 0; isProcessed = true {
		stackItem := componentStack[len(componentStack)-1]
		for _, identifier := range stackItem.children {
			if visited[identifier] {
				continue
			}
			visited[identifier] = true
			componentStack = append(componentStack, componentMap[identifier])
			isProcessed = false
		}
		if !isProcessed {
			continue
		}
		// We reduce the stack in any case, but the component will be added to the result slice only if it is not there.
		componentStack = componentStack[:len(componentStack)-1]
		if !processed[stackItem.identifier] {
			processed[stackItem.identifier] = true
			result = append(result, stackItem.component)
		}
	}
	return result
}

type stackItem struct {
	identifier string
	component  componego.Component
	parent     *stackItem
	children   []string
	position   int
}
