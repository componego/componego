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

package application

import (
	"errors"

	"github.com/componego/componego"
)

type Factory interface {
	SetApplicationName(name string)
	SetApplicationComponents(components func() ([]componego.Component, error))
	SetApplicationDependencies(dependencies func() ([]componego.Dependency, error))
	SetApplicationConfigInit(configInit func(appMode componego.ApplicationMode, options any) (map[string]any, error))
	SetApplicationErrorHandler(errorHandler func(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error)
	SetApplicationAction(action func(env componego.Environment, options any) (int, error))
	Build() componego.Application
}

type factory struct {
	name         string
	components   func() ([]componego.Component, error)
	dependencies func() ([]componego.Dependency, error)
	configInit   func(appMode componego.ApplicationMode, options any) (map[string]any, error)
	errorHandler func(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error
	action       func(env componego.Environment, options any) (int, error)
}

func NewFactory(name string) Factory {
	return &factory{
		name: name,
	}
}

// SetApplicationName belongs to interface Factory.
func (f *factory) SetApplicationName(name string) {
	f.name = name
}

// SetApplicationComponents belongs to interface Factory.
func (f *factory) SetApplicationComponents(components func() ([]componego.Component, error)) {
	f.components = components
}

// SetApplicationDependencies belongs to interface Factory.
func (f *factory) SetApplicationDependencies(dependencies func() ([]componego.Dependency, error)) {
	f.dependencies = dependencies
}

// SetApplicationConfigInit belongs to interface Factory.
func (f *factory) SetApplicationConfigInit(configInit func(appMode componego.ApplicationMode, options any) (map[string]any, error)) {
	f.configInit = configInit
}

// SetApplicationErrorHandler belongs to interface Factory.
func (f *factory) SetApplicationErrorHandler(errorHandler func(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error) {
	f.errorHandler = errorHandler
}

// SetApplicationAction belongs to interface Factory.
func (f *factory) SetApplicationAction(action func(env componego.Environment, options any) (int, error)) {
	f.action = action
}

// Build belongs to interface Factory.
func (f *factory) Build() componego.Application {
	return &QuickApplication{
		Name:         f.name,
		Components:   f.components,
		Dependencies: f.dependencies,
		ConfigInit:   f.configInit,
		ErrorHandler: f.errorHandler,
		Action:       f.action,
	}
}

type QuickApplication struct {
	Name         string
	Components   func() ([]componego.Component, error)
	Dependencies func() ([]componego.Dependency, error)
	ConfigInit   func(appMode componego.ApplicationMode, options any) (map[string]any, error)
	ErrorHandler func(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error
	Action       func(env componego.Environment, options any) (int, error)
}

// ApplicationName belongs to interface componego.Application.
func (q *QuickApplication) ApplicationName() string {
	return q.Name
}

// ApplicationComponents belongs to interface componego.ApplicationComponents.
func (q *QuickApplication) ApplicationComponents() ([]componego.Component, error) {
	if q.Components == nil {
		return nil, nil
	}
	return q.Components()
}

// ApplicationDependencies belongs to interface componego.ApplicationDependencies.
func (q *QuickApplication) ApplicationDependencies() ([]componego.Dependency, error) {
	if q.Dependencies == nil {
		return nil, nil
	}
	return q.Dependencies()
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (q *QuickApplication) ApplicationConfigInit(appMode componego.ApplicationMode, options any) (map[string]any, error) {
	if q.ConfigInit == nil {
		return nil, nil
	}
	return q.ConfigInit(appMode, options)
}

// ApplicationErrorHandler belongs to interface componego.ApplicationErrorHandler.
func (q *QuickApplication) ApplicationErrorHandler(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error {
	if q.ErrorHandler == nil {
		return err
	}
	return q.ErrorHandler(err, appIO, appMode)
}

// ApplicationAction belongs to interface componego.Application.
func (q *QuickApplication) ApplicationAction(env componego.Environment, options any) (int, error) {
	if q.Action == nil {
		return ExitWrapper(errors.New("there is no action in application"))
	}
	return q.Action(env, options)
}

var (
	_ componego.Application             = (*QuickApplication)(nil)
	_ componego.ApplicationComponents   = (*QuickApplication)(nil)
	_ componego.ApplicationDependencies = (*QuickApplication)(nil)
	_ componego.ApplicationConfigInit   = (*QuickApplication)(nil)
	_ componego.ApplicationErrorHandler = (*QuickApplication)(nil)
)
