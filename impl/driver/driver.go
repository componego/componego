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
	"errors"
	"reflect"
	"runtime"

	"github.com/componego/componego"
	"github.com/componego/componego/libs/debug"
	"github.com/componego/componego/libs/xerrors"
)

var (
	// ErrPanic is an error thrown when a panic is received.
	ErrPanic = xerrors.New("global panic")
	// ErrUnknownPanic is an error thrown when a panic of an unknown type is received.
	ErrUnknownPanic = ErrPanic.WithMessage("unknown panic")
)

// Driver is an interface that describes the driver for starting and controlling the application.
type Driver interface {
	// RunApplication launches the application.
	RunApplication(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (int, error)
	// CreateEnvironment creates a new environment for the application.
	CreateEnvironment(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (componego.Environment, error)
}

type driver struct {
	options *Options
}

// New is an application driver constructor.
func New(options *Options) Driver {
	return &driver{
		// It adds options if any of the options are missing.
		options: Configure(options),
	}
}

// RunApplication launches the application.
func (d *driver) RunApplication(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (exitCode int, err error) {
	cancelableCtx, cancelCtx := context.WithCancel(ctx)
	defer func() {
		cancelCtx()       // All contexts that were created from the main context will be closed on exit.
		runtime.Gosched() // We switch the runtime so that waiting goroutines can complete their work.
		// It converts the panic into an error with the current stack as an error option.
		if err = ErrorRecoveryOnStop(recover(), err); err == nil {
			return
		}
		// We handle all application errors that were not previously handled.
		if app, ok := app.(componego.ApplicationErrorHandler); ok {
			err = app.ApplicationErrorHandler(err, d.options.AppIO, appMode)
		}
		// If the application exit code is successful, but at some stage an error occurred,
		// then we will change the code to unsuccessful.
		if exitCode == componego.SuccessExitCode {
			exitCode = componego.ErrorExitCode
		}
	}()
	if env, err := d.CreateEnvironment(cancelableCtx, app, appMode); err != nil {
		return componego.ErrorExitCode, err
	} else { //nolint:revive
		return d.runInsideEnvironment(env)
	}
}

// CreateEnvironment creates a new environment for the application.
func (d *driver) CreateEnvironment(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (env componego.Environment, err error) {
	defer func() {
		// This function should not return panic.
		// Errors of goroutines running inside the environment must be processed in the places where they are launched.
		err = ErrorRecoveryOnStop(recover(), err)
	}()
	// One driver can run multiple applications. Therefore, we don't use managers directly.
	// Each application will have its own instances of managers and environment.
	// This way we guarantee that the values of the variables will not overlap.
	configProvider, configInitializer := d.options.ConfigProviderFactory()
	componentProvider, componentsInitializer := d.options.ComponentProviderFactory()
	dependencyInvoker, dependenciesInitializer := d.options.DependencyInvokerFactory()
	// Pay attention to the order in which the functions are run.
	env = d.options.EnvironmentFactory(
		ctx, app, d.options.AppIO, appMode, configProvider, componentProvider, dependencyInvoker,
	)
	if err = configInitializer(env); err != nil {
		return nil, err
	}
	if err = componentsInitializer(env); err != nil {
		return nil, err
	}
	if err = dependenciesInitializer(env); err != nil {
		return nil, err
	}
	return env, nil
}

func (d *driver) runInsideEnvironment(env componego.Environment) (exitCode int, err error) {
	for _, component := range env.Components() {
		// The order in which components are called depends on the dependencies between the components.
		// Therefore, it is very important to indicate which components your component depends on.
		if component, ok := component.(componego.ComponentInit); ok {
			if err = component.ComponentInit(env); err != nil {
				return componego.ErrorExitCode, err
			}
		}
		if component, ok := component.(componego.ComponentStop); ok {
			// Stopping a component is guaranteed to occur in the reverse order of component initialization.
			// noinspection ALL
			defer func(component componego.ComponentStop) {
				err = ErrorRecoveryOnStop(recover(), err) // We catch the panic that may occur.
				err = component.ComponentStop(env, err)   // It can handle this error somehow or/and return it to work.
			}(component) // We support compatibility with older versions of the language.
		}
	}
	return env.Application().ApplicationAction(env, d.options.Args)
}

// ErrorRecoveryOnStop returns an error after panic recovery.
func ErrorRecoveryOnStop(recover any, prevErr error) (newErr error) {
	if recover == nil {
		return prevErr
	}
	errOptions := []xerrors.Option{
		xerrors.NewOption("componego:driver:panic:stack", debug.GetStackTrace()),
		xerrors.NewOption("componego:driver:panic:recover", recover),
	}
	// noinspection ALL
	if err, ok := recover.(error); ok {
		newErr = ErrPanic.WithError(err, errOptions...)
	} else if message, ok := recover.(string); ok {
		newErr = ErrPanic.WithMessage(message, errOptions...)
	} else if reflectValue := reflect.ValueOf(recover); reflectValue.Kind() == reflect.String {
		// Handle the remaining type set of ~string.
		newErr = ErrPanic.WithMessage(reflectValue.String(), errOptions...)
	} else {
		newErr = ErrUnknownPanic.WithOptions(errOptions...)
	}
	if prevErr == nil {
		return newErr
	}
	return errors.Join(prevErr, newErr)
}
