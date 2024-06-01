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

import (
	"io"
)

// List of default application modes.
const (
	// ProductionMode is a constant for running the application in production mode.
	ProductionMode ApplicationMode = 0
	// DeveloperMode is a constant for running the application in developer mode.
	DeveloperMode ApplicationMode = 1
	// TestMode is a constant for running the application in test mode.
	TestMode ApplicationMode = 2
)

// List of default application exit codes.
const (
	// SuccessExitCode is the exit code that the application should return
	// if there were no errors and the application executed successfully.
	SuccessExitCode int = 0
	// ErrorExitCode is the exit code that the application should return if the application fails.
	// It can also be any number that is greater than 0 and less than 256.
	ErrorExitCode int = 1
)

// ApplicationMode is application mode type.
type ApplicationMode int

// Application is a general interface that application should describe.
type Application interface {
	// ApplicationName returns the name of the application.
	ApplicationName() string
	// ApplicationAction describes the main action of the current application.
	// This function is called last when the application is fully initialized
	ApplicationAction(env Environment, args []string) (int, error)
}

// ApplicationComponents is an interface that describes the components of the application.
type ApplicationComponents interface {
	// Application belongs to the application.
	Application
	// ApplicationComponents return a list of components that the application depends on.
	// They will be sorted based on the dependent components.
	ApplicationComponents() ([]Component, error)
}

// ApplicationDependencies is an interface that describes the dependencies of the application.
type ApplicationDependencies interface {
	// Application belongs to the application.
	Application
	// ApplicationDependencies returns a list of dependencies that the application provides.
	// This function must return an array of objects or functions that returns objects.
	ApplicationDependencies() ([]Dependency, error)
}

// ApplicationConfigInit is an interface that describes the process of obtaining configuration for the application.
type ApplicationConfigInit interface {
	// Application belongs to the application.
	Application
	// ApplicationConfigInit returns configuration for all application entities.
	// This function is called only once.
	ApplicationConfigInit(appMode ApplicationMode) (map[string]any, error)
}

// ApplicationErrorHandler is an interface that describes a function that catches all application errors.
type ApplicationErrorHandler interface {
	// Application belongs to the application.
	Application
	// ApplicationErrorHandler catches all application errors that are not handled in the main application code.
	ApplicationErrorHandler(err error, appIO ApplicationIO, appMode ApplicationMode) error
}

// ApplicationIO is an interface that describes input and output for the application environment.
type ApplicationIO interface {
	// InputReader returns a pointer to read data into the application.
	InputReader() io.Reader
	// OutputWriter returns a pointer to output data from the application.
	OutputWriter() io.Writer
	// ErrorOutputWriter This function returns a pointer to output errors from the application.
	ErrorOutputWriter() io.Writer
}
