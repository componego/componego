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

package environment

import (
	"context"
	"sync"

	"github.com/componego/componego"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/xerrors"
)

var (
	ErrInvalidParentContext   = xerrors.New("new content is not created based on previous context", "E0210")
	ErrNoEnvironmentInContext = xerrors.New("there is no environment in the context", "E0220")
)

type contextKey struct{}

type environment struct {
	mutex             sync.Mutex
	context           context.Context
	application       componego.Application
	applicationIO     componego.ApplicationIO
	applicationMode   componego.ApplicationMode
	configProvider    componego.ConfigProvider
	componentProvider componego.ComponentProvider
	dependencyInvoker componego.DependencyInvoker
}

// New is a constructor that creates environment for the application.
func New(
	ctx context.Context,
	application componego.Application,
	applicationIO componego.ApplicationIO,
	applicationMode componego.ApplicationMode,
	configProvider componego.ConfigProvider,
	componentProvider componego.ComponentProvider,
	dependencyInvoker componego.DependencyInvoker,
) componego.Environment {
	// Ignore the fact that the function takes too many arguments. In this case, this is the best way at this stage.
	env := &environment{
		mutex:             sync.Mutex{},
		application:       application,
		applicationIO:     applicationIO,
		applicationMode:   applicationMode,
		configProvider:    configProvider,
		componentProvider: componentProvider,
		dependencyInvoker: dependencyInvoker,
	}
	env.context = context.WithValue(ctx, contextKey{}, env)
	return env
}

// GetContext returns a current application context.
func (e *environment) GetContext() context.Context {
	e.mutex.Lock()
	ctx := e.context
	e.mutex.Unlock()
	return ctx
}

// SetContext sets a new application context.
func (e *environment) SetContext(ctx context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// We make sure that the new context was created based on the previous context.
	if utils.IsParentContext(e.context, ctx) {
		e.context = ctx
		return nil
	}
	return ErrInvalidParentContext
}

// Application returns a current application object.
func (e *environment) Application() componego.Application {
	return e.application
}

// ApplicationIO returns an object for getting application input and output.
func (e *environment) ApplicationIO() componego.ApplicationIO {
	return e.applicationIO
}

// ApplicationMode returns the mode in which the application is started.
// This is the value that you pass to the runner during application initialization.
func (e *environment) ApplicationMode() componego.ApplicationMode {
	return e.applicationMode
}

// ConfigProvider returns an object for getting config.
func (e *environment) ConfigProvider() componego.ConfigProvider {
	return e.configProvider
}

// Components returns a sorted list of active application components.
func (e *environment) Components() []componego.Component {
	return e.componentProvider.Components()
}

// DependencyInvoker returns an object to invoke dependencies.
func (e *environment) DependencyInvoker() componego.DependencyInvoker {
	return e.dependencyInvoker
}

func GetEnvironment(ctx context.Context) (componego.Environment, error) {
	if env, ok := ctx.Value(contextKey{}).(componego.Environment); ok {
		return env, nil
	}
	return nil, ErrNoEnvironmentInContext
}

func GetEnvironmentOrPanic(ctx context.Context) componego.Environment {
	env, err := GetEnvironment(ctx)
	if err != nil {
		panic(err)
	}
	return env
}
