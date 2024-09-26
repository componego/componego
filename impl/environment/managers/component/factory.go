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

package component

import (
	"github.com/componego/componego"
)

type Factory interface {
	SetComponentIdentifier(identifier string)
	SetComponentVersion(version string)
	SetComponentComponents(components func() ([]componego.Component, error))
	SetComponentDependencies(dependencies func() ([]componego.Dependency, error))
	SetComponentInit(init func(env componego.Environment) error)
	SetComponentStop(stop func(env componego.Environment, prevErr error) error)
	Build() componego.Component
}

type factory struct {
	identifier   string
	version      string
	components   func() ([]componego.Component, error)
	dependencies func() ([]componego.Dependency, error)
	init         func(env componego.Environment) error
	stop         func(env componego.Environment, prevErr error) error
}

func NewFactory(identifier string, version string) Factory {
	return &factory{
		identifier: identifier,
		version:    version,
	}
}

// SetComponentIdentifier belongs to interface Factory.
func (f *factory) SetComponentIdentifier(identifier string) {
	f.identifier = identifier
}

// SetComponentVersion belongs to interface Factory.
func (f *factory) SetComponentVersion(version string) {
	f.version = version
}

// SetComponentComponents belongs to interface Factory.
func (f *factory) SetComponentComponents(components func() ([]componego.Component, error)) {
	f.components = components
}

// SetComponentDependencies belongs to interface Factory.
func (f *factory) SetComponentDependencies(dependencies func() ([]componego.Dependency, error)) {
	f.dependencies = dependencies
}

// SetComponentInit belongs to interface Factory.
func (f *factory) SetComponentInit(init func(env componego.Environment) error) {
	f.init = init
}

// SetComponentStop belongs to interface Factory.
func (f *factory) SetComponentStop(stop func(env componego.Environment, prevErr error) error) {
	f.stop = stop
}

// Build belongs to interface Factory.
func (f *factory) Build() componego.Component {
	return &QuickComponent{
		Identifier:   f.identifier,
		Version:      f.version,
		Components:   f.components,
		Dependencies: f.dependencies,
		Init:         f.init,
		Stop:         f.stop,
	}
}

type QuickComponent struct {
	Identifier   string
	Version      string
	Components   func() ([]componego.Component, error)
	Dependencies func() ([]componego.Dependency, error)
	Init         func(env componego.Environment) error
	Stop         func(env componego.Environment, prevErr error) error
}

// ComponentIdentifier belongs to interface componego.Component.
func (q *QuickComponent) ComponentIdentifier() string {
	return q.Identifier
}

// ComponentVersion belongs to interface componego.Component.
func (q *QuickComponent) ComponentVersion() string {
	return q.Version
}

// ComponentComponents belongs to interface componego.ComponentComponents.
func (q *QuickComponent) ComponentComponents() ([]componego.Component, error) {
	if q.Components == nil {
		return nil, nil
	}
	return q.Components()
}

// ComponentDependencies belongs to interface componego.ComponentDependencies.
func (q *QuickComponent) ComponentDependencies() ([]componego.Dependency, error) {
	if q.Dependencies == nil {
		return nil, nil
	}
	return q.Dependencies()
}

// ComponentInit belongs to interface componego.ComponentInit.
func (q *QuickComponent) ComponentInit(env componego.Environment) error {
	if q.Init == nil {
		return nil
	}
	return q.Init(env)
}

// ComponentStop belongs to interface componego.ComponentStop.
func (q *QuickComponent) ComponentStop(env componego.Environment, prevErr error) error {
	if q.Stop == nil {
		return prevErr
	}
	return q.Stop(env, prevErr)
}

var (
	_ componego.Component             = (*QuickComponent)(nil)
	_ componego.ComponentComponents   = (*QuickComponent)(nil)
	_ componego.ComponentDependencies = (*QuickComponent)(nil)
	_ componego.ComponentInit         = (*QuickComponent)(nil)
	_ componego.ComponentStop         = (*QuickComponent)(nil)
)
