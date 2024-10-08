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

package driver

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"sync"

	"github.com/componego/componego"
	"github.com/componego/componego/libs/debug"
	"github.com/componego/componego/libs/xerrors"
)

var (
	// ErrPanic is an error thrown when a panic is received.
	ErrPanic = xerrors.New("global panic", "E0100")
	// ErrUnknownPanic is an error thrown when a panic of an unknown type is received.
	ErrUnknownPanic = ErrPanic.WithMessage("unknown panic", "E0101")
)

// Driver is an interface that describes the driver for starting and controlling the application.
type Driver interface {
	// RunApplication launches the application.
	RunApplication(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (int, error)
	// CreateEnvironment creates a new environment for the application.
	CreateEnvironment(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (componego.Environment, canceller, error)
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
	defer func() {
		runtime.Gosched() // We switch the runtime so that waiting goroutines can complete their work.
		if err == nil {   // All panics are caught in previous functions.
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
	env, cancelEnv, errEnv := d.CreateEnvironment(ctx, app, appMode)
	if errEnv != nil {
		return componego.ErrorExitCode, errEnv
	}
	defer func() {
		// It converts the panic into an error with the current stack as an error option.
		err = ErrorRecoveryOnStop(recover(), err)
		err = errors.Join(err, cancelEnv())
	}()
	return d.runInsideEnvironment(env)
}

// CreateEnvironment creates a new environment for the application.
func (d *driver) CreateEnvironment(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) (env componego.Environment, cancelEnv canceller, err error) {
	defer func() {
		// This function should not return panic.
		// Errors of goroutines running inside the environment must be processed in the places where they are launched.
		if err = ErrorRecoveryOnStop(recover(), err); err != nil {
			err = errors.Join(err, cancelEnv())
			cancelEnv = nil
		}
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
	// We notify environment managers about termination after a normal application shutdown or when an error occurs.
	cancels := make([]canceller, 3)
	cancelEnv = sync.OnceValue[error](joinCancels(cancels))
	if cancels[0], err = configInitializer(env, d.options.Additional); err != nil {
		return nil, cancelEnv, err
	}
	if cancels[1], err = componentsInitializer(env, d.options.Additional); err != nil {
		return nil, cancelEnv, err
	}
	if cancels[2], err = dependenciesInitializer(env, d.options.Additional); err != nil {
		return nil, cancelEnv, err
	}
	return env, cancelEnv, nil
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
				runtime.Gosched()                         // We switch the runtime so that the waiting goroutines can stop their work.
				err = ErrorRecoveryOnStop(recover(), err) // We catch the panic that may occur.
				err = component.ComponentStop(env, err)   // It can handle this error somehow or/and return it to work.
			}(component) // We support compatibility with older versions of the language.
		}
	}
	return env.Application().ApplicationAction(env, d.options.Additional)
}

// ErrorRecoveryOnStop returns an error after panic recovery.
func ErrorRecoveryOnStop(recover any, prevErr error) (newErr error) {
	if recover == nil {
		return prevErr
	}
	errOptions := []xerrors.Option{
		xerrors.NewOption("componego:driver:panic:stack", debug.GetStackTrace(1)),
		xerrors.NewOption("componego:driver:panic:recover", recover),
	}
	// noinspection ALL
	if err, ok := recover.(error); ok {
		newErr = ErrPanic.WithError(err, "E0102", errOptions...)
	} else if message, ok := recover.(string); ok {
		newErr = ErrPanic.WithMessage(message, "E0103", errOptions...)
	} else if reflectValue := reflect.ValueOf(recover); reflectValue.Kind() == reflect.String {
		// Handle the remaining type set of ~string.
		newErr = ErrPanic.WithMessage(reflectValue.String(), "E0104", errOptions...)
	} else {
		newErr = ErrUnknownPanic.WithOptions("E0105", errOptions...)
	}
	if prevErr == nil {
		return newErr
	}
	return errors.Join(prevErr, newErr)
}

func joinCancels(cancels []canceller) canceller {
	return func() (err error) {
		defer func() {
			// This catches the error in the last called function.
			err = ErrorRecoveryOnStop(recover(), err)
		}()
		// All functions will be called starting from the last function in the list.
		// if an error occurs in any of the functions, it will be merged with other errors.
		for _, cancel := range cancels {
			if cancel == nil {
				continue
			}
			// noinspection ALL
			defer func(cancel canceller) {
				err = ErrorRecoveryOnStop(recover(), err)
				err = errors.Join(err, cancel())
			}(cancel)
		}
		return err
	}
}
