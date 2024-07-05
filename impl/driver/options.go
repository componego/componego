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

package driver

import (
	"context"
	"os"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/environment"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/impl/environment/managers/config"
	"github.com/componego/componego/impl/environment/managers/dependency"
	"github.com/componego/componego/impl/environment/managers/dependency/container"
	"github.com/componego/componego/internal/system"
)

type initializer = func(env componego.Environment, options any) error

type Options struct {
	ConfigProviderFactory    func() (componego.ConfigProvider, initializer)
	ComponentProviderFactory func() (componego.ComponentProvider, initializer)
	DependencyInvokerFactory func() (componego.DependencyInvoker, initializer)
	EnvironmentFactory       func(
		context context.Context,
		application componego.Application,
		applicationIO componego.ApplicationIO,
		applicationMode componego.ApplicationMode,
		configProvider componego.ConfigProvider,
		componentProvider componego.ComponentProvider,
		dependencyInvoker componego.DependencyInvoker,
	) componego.Environment
	AppIO      componego.ApplicationIO
	Additional any
}

func Configure(options *Options) *Options {
	if options == nil {
		options = &Options{}
	}
	if options.ConfigProviderFactory == nil {
		options.ConfigProviderFactory = newConfigFactory
	}
	if options.ComponentProviderFactory == nil {
		options.ComponentProviderFactory = newComponentProviderFactory
	}
	if options.DependencyInvokerFactory == nil {
		options.DependencyInvokerFactory = newDependencyInvokerFactory
	}
	if options.EnvironmentFactory == nil {
		options.EnvironmentFactory = environment.New
	}
	if options.AppIO == nil {
		options.AppIO = application.NewIO(system.Stdin, system.Stdout, system.Stderr)
	}
	if options.Additional == nil {
		// This variable can contain any data depending on how the application is started.
		// By default, these are command line arguments.
		options.Additional = os.Args
	}
	return options
}

func newComponentProviderFactory() (componego.ComponentProvider, initializer) {
	manager, initializer := component.NewManager()
	return manager, func(env componego.Environment, _ any) error {
		components, err := component.ExtractComponents(env.Application())
		if err != nil {
			return err
		}
		return initializer(components)
	}
}

func newDependencyInvokerFactory() (componego.DependencyInvoker, initializer) {
	manager, initializer := dependency.NewManager()
	return manager, func(env componego.Environment, _ any) error {
		dependencies, err := dependency.ExtractDependencies(env)
		if err != nil {
			return err
		}
		containerInstance, containerInitializer := container.New(len(dependencies))
		// There may be a recursive call to the container through the dependency manager
		// during the initialization of dependencies inside the container.
		if err = initializer(containerInstance); err != nil {
			return err
		}
		return containerInitializer(dependencies)
	}
}

func newConfigFactory() (componego.ConfigProvider, initializer) {
	manager, initializer := config.NewManager()
	return manager, func(env componego.Environment, options any) error {
		parsedConfig, err := config.ParseConfig(env, options)
		if err != nil {
			return err
		}
		return initializer(env, parsedConfig)
	}
}
