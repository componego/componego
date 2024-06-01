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

package componego

// Component is an interface that describes the component.
type Component interface {
	// ComponentIdentifier returns the component ID.
	// If the identifier in several components is the same, then the last component will be used.
	// This can be used to overwrite components.
	ComponentIdentifier() string
	// ComponentVersion returns the component version.
	ComponentVersion() string
}

// ComponentComponents is an interface that describes which components the current component depends on.
type ComponentComponents interface {
	// Component belongs to the component.
	Component
	// ComponentComponents returns list of components that the component depends on.
	// They will be sorted based on the dependent components.
	ComponentComponents() ([]Component, error)
}

// ComponentDependencies is an interface that describes the dependencies of the component.
type ComponentDependencies interface {
	// Component belongs to the component.
	Component
	// ComponentDependencies returns a list of dependencies that the component provides
	// This function must return an array of objects or functions that returns objects.
	ComponentDependencies() ([]Dependency, error)
}

// ComponentInit is an interface that describes the component initialization.
type ComponentInit interface {
	// Component belongs to the component.
	Component
	// ComponentInit is called during component initialization.
	ComponentInit(env Environment) error
}

// ComponentStop is an interface that describes stopping the component.
type ComponentStop interface {
	// Component belongs to the component.
	Component
	// ComponentStop is called when the component stops.
	// You can handle previous error (return a new or old error).
	ComponentStop(env Environment, prevErr error) error
}

// ComponentProvider is an interface that describes a list of active application components.
// These components are sorted in order of dependencies between components.
type ComponentProvider interface {
	// Components returns a list of sorted application components.
	Components() []Component
}
