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

package componego

import (
	"context"
)

// Environment is an interface that describes the application environment.
type Environment interface {
	// GetContext returns a current application context.
	GetContext() context.Context
	// SetContext sets a new application context.
	// The new context must inherit from the previous context.
	SetContext(ctx context.Context) error
	// Application returns a current application object.
	Application() Application
	// ApplicationIO returns an object for getting application input and output.
	ApplicationIO() ApplicationIO
	// ApplicationMode returns the mode in which the application is started.
	ApplicationMode() ApplicationMode
	// ConfigProvider returns an object for getting config.
	ConfigProvider() ConfigProvider
	// Components returns a sorted list of active application components.
	Components() []Component
	// DependencyInvoker returns an object to invoke dependencies.
	DependencyInvoker() DependencyInvoker
}
