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

package dependency

import (
	"context"
	"reflect"
	"strings"
	"unsafe"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency/container"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/xerrors"
)

var (
	ErrDependencyManager   = xerrors.New("error inside dependency manager")
	ErrExtractDependencies = ErrDependencyManager.WithMessage("an error occurred while receiving dependencies")
	ErrNilArgument         = ErrDependencyManager.WithMessage("nil argument")
	ErrNotFunction         = ErrDependencyManager.WithMessage("argument is not a function and cannot be used as a constructor for dependency injection")
	ErrVariadicFunction    = ErrDependencyManager.WithMessage("function has a variable number of arguments and cannot be used as a constructor for dependency injection")
	ErrNotAllowedTarget    = ErrDependencyManager.WithMessage("target is not allowed for dependency injection")
)

type manager struct {
	container container.Container
}

func NewManager() (componego.DependencyInvoker, func(container.Container) error) {
	m := &manager{}
	return m, m.initialize
}

func (m *manager) Invoke(function any) (any, error) {
	if function == nil {
		return nil, ErrNilArgument
	}
	reflectType := reflect.TypeOf(function)
	if reflectType.Kind() != reflect.Func {
		return nil, ErrNotFunction.WithOptions(
			xerrors.NewOption("componego:dependency:function", reflectType),
		)
	}
	if reflectType.IsVariadic() {
		return nil, ErrVariadicFunction.WithOptions(
			xerrors.NewOption("componego:dependency:function", reflectType),
		)
	}
	numIn := reflectType.NumIn()
	dependencies := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		value, err := m.container.GetValue(reflectType.In(i))
		if err != nil {
			return nil, ErrDependencyManager.WithError(err,
				xerrors.NewOption("componego:dependency:function", reflectType),
				xerrors.NewOption("componego:dependency:argument", reflectType.In(i)),
			)
		}
		dependencies[i] = value
	}
	output := reflect.ValueOf(function).Call(dependencies)
	if len(output) == 0 {
		return nil, nil
	}
	outputLen := len(output)
	if last := output[outputLen-1]; utils.IsErrorType(last.Type()) {
		if errInstance := last.Interface(); errInstance != nil {
			return nil, ErrDependencyManager.WithError(errInstance.(error),
				xerrors.NewOption("componego:dependency:function", reflectType),
			)
		}
		outputLen--
	}
	if outputLen == 0 {
		return nil, nil
	}
	return output[0].Interface(), nil
}

func (m *manager) Populate(target any) error {
	if target == nil {
		return ErrNilArgument
	}
	reflectType := reflect.TypeOf(target).Elem()
	if reflectType.Kind() != reflect.Interface &&
		(reflectType.Kind() != reflect.Pointer || reflectType.Elem().Kind() != reflect.Struct) {
		return ErrNotAllowedTarget.WithOptions(
			xerrors.NewOption("componego:dependency:target", reflectType),
		)
	}
	value, err := m.container.GetValue(reflectType)
	if err != nil {
		return ErrDependencyManager.WithError(err,
			xerrors.NewOption("componego:dependency:target", reflectType),
		)
	}
	// We are sure that this is a pointer because the target was validated.
	reflect.ValueOf(target).Elem().Set(value)
	return nil
}

func (m *manager) PopulateFields(target any) error {
	if target == nil {
		return ErrNilArgument
	}
	reflectType := reflect.TypeOf(target)
	if reflectType.Kind() != reflect.Pointer || reflectType.Elem().Kind() != reflect.Struct {
		return ErrNotAllowedTarget.WithOptions(
			xerrors.NewOption("componego:dependency:target", reflectType),
		)
	}
	reflectType = reflectType.Elem()
	reflectValue := reflect.ValueOf(target).Elem()
	numField := reflectType.NumField()
	for i := 0; i < numField; i++ {
		field := reflectType.Field(i)
		if field.Tag != `componego:"inject"` { // minor optimization.
			tag, ok := field.Tag.Lookup("componego")
			if !ok || (tag != "inject" && !utils.Contains(strings.Split(tag, ","), "inject")) {
				continue
			}
		}
		value, err := m.container.GetValue(field.Type)
		if err != nil {
			return ErrDependencyManager.WithError(err,
				xerrors.NewOption("componego:dependency:target", reflect.TypeOf(target)), // original type.
				xerrors.NewOption("componego:dependency:fieldName", field.Name),
			)
		}
		if field := reflectValue.Field(i); field.CanSet() {
			field.Set(value)
		} else {
			// Support for non-exported fields.
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(value) // #nosec G103
		}
	}
	return nil
}

func (m *manager) initialize(container container.Container) error {
	m.container = container
	return nil
}

// ExtractDependencies returns a list of dependencies from the application and components.
// This is a raw list without any transformations.
func ExtractDependencies(env componego.Environment) ([]componego.Dependency, error) {
	components := env.Components()
	allDependencies := make([][]componego.Dependency, 0, len(components)+1)
	countDependencies := 0
	for _, component := range components {
		if component, ok := component.(componego.ComponentDependencies); !ok {
			continue
		} else if dependencies, err := component.ComponentDependencies(); err != nil {
			return nil, ErrExtractDependencies.WithError(err,
				xerrors.NewOption("componego:dependency:component", component),
			)
		} else if len(dependencies) > 0 {
			allDependencies = append(allDependencies, dependencies)
			countDependencies += len(dependencies)
		}
	}
	if app, ok := env.Application().(componego.ApplicationDependencies); ok {
		if dependencies, err := app.ApplicationDependencies(); err != nil {
			return nil, ErrExtractDependencies.WithError(err)
		} else if len(dependencies) > 0 {
			allDependencies = append(allDependencies, dependencies)
			countDependencies += len(dependencies)
		}
	}
	defaultDependencies := getDefaultDependencies(env)
	dependencies := make([]componego.Dependency, 0, countDependencies+len(defaultDependencies))
	for _, list := range allDependencies {
		dependencies = append(dependencies, list...)
	}
	// Adds dependencies that will be present in any application and cannot be overwritten (because they are added at the end).
	dependencies = append(dependencies, defaultDependencies...)
	return dependencies, nil
}

// getDefaultDependencies returns dependencies that will be present in any application.
func getDefaultDependencies(env componego.Environment) []componego.Dependency {
	return []componego.Dependency{
		func() componego.Environment {
			return env
		},
		func() context.Context {
			// Context cannot be provided as a dependency.
			return nil
		},
		env.Application,
		env.ApplicationIO,
		env.ConfigProvider,
		env.DependencyInvoker,
	}
}
