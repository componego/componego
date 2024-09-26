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

package runner

import (
	"context"
	"errors"
	"io"
	"runtime"
	"sync"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/internal/testing"
)

type TestOptions struct {
	Driver             driver.Driver
	EnvironmentFactory func(app componego.Application, driver driver.Driver) (componego.Environment, func() error, error)
	OnApplicationStop  func(app componego.Application, returnErr error) error
	OnComponentInit    func(env componego.Environment, component componego.ComponentInit, returnErr error)
	OnComponentStop    func(env componego.Environment, component componego.ComponentStop, returnErr error, previousErr error)
}

func (t *TestOptions) createEnvironment(app componego.Application) (componego.Environment, func() error, error) {
	if t.EnvironmentFactory != nil {
		// If you want to change the base content or launch mode, you can do it here.
		return t.EnvironmentFactory(app, t.getDriver())
	}
	return t.getDriver().CreateEnvironment(context.Background(), app, componego.TestMode)
}

func (t *TestOptions) getDriver() driver.Driver {
	if t.Driver != nil {
		// Multiple applications can share a driver.
		return t.Driver
	}
	return driver.New(&driver.Options{
		// Ignore output. We don't need output during the test.
		AppIO: application.NewIO(nil, io.Discard, io.Discard),
	})
}

func (t *TestOptions) onComponentInit(env componego.Environment, component componego.ComponentInit, returnErr error) error {
	if t.OnComponentInit != nil {
		// This is the hook where you can check the returned error after a component is initialized. .
		t.OnComponentInit(env, component, returnErr)
	}
	return returnErr // The hook does not change the returned error.
}

func (t *TestOptions) onComponentStop(env componego.Environment, component componego.ComponentStop, returnErr error, previousErr error) error {
	if t.OnComponentStop != nil {
		// This is the hook where you can check the returned error after a component is stopped.
		t.OnComponentStop(env, component, returnErr, previousErr)
	}
	return returnErr // The hook does not change the returned error.
}

func (t *TestOptions) onApplicationStop(app componego.Application, returnErr error) error {
	if t.OnApplicationStop == nil {
		return returnErr
	}
	// You may not return an error from the hook if the application should fail, but you don't want the test to fail.
	return t.OnApplicationStop(app, returnErr)
}

func CreateTestEnvironment(t testing.T, app componego.Application, options *TestOptions) (componego.Environment, func()) {
	if options == nil {
		options = &TestOptions{}
	}
	env, cancelEnv, err := options.createEnvironment(app)
	if err != nil {
		t.Errorf("error when creating an environment for the application: %s", err)
		t.FailNow()
		return nil, nil
	}
	components := env.Components()
	stopComponents := make([]componego.ComponentStop, 0, len(components))
	// All processes can be stopped only once.
	cancelAll := sync.OnceFunc(func() {
		defer func() {
			// This function should not return panic.
			err = errors.Join(driver.ErrorRecoveryOnStop(recover(), err), cancelEnv())
			if err = options.onApplicationStop(app, err); err != nil {
				// If the application fails, that error will be caught in the hook.
				// You may not return an error from the hook if the application should fail, but you don't want the test to fail.
				t.Errorf("error when stopping the application: %s", err)
				t.FailNow()
			}
		}()
		for _, component := range stopComponents {
			// All components will be stopped in a reverse of initialization order.
			// noinspection ALL
			defer func(component componego.ComponentStop) {
				// We switch the runtime so that waiting goroutines can complete their work.
				runtime.Gosched()
				// We catch the panic that could have been in the previous components.
				err = driver.ErrorRecoveryOnStop(recover(), err)
				// Each call can be checked in the hook specified in the options.
				// For example, you can ensure the component's stop function has returned and processes previous errors correctly.
				err = options.onComponentStop(env, component, component.ComponentStop(env, err), err)
			}(component)
		}
	})
	// All processes will be cancelled even if the cancel function is not called manually.
	t.Cleanup(cancelAll)
	defer func() {
		if err = driver.ErrorRecoveryOnStop(recover(), err); err != nil {
			cancelAll() // If an error occurred during component initialization.
		}
	}()
	for _, component := range components {
		if component, ok := component.(componego.ComponentInit); ok {
			if err = options.onComponentInit(env, component, component.ComponentInit(env)); err != nil {
				return nil, nil
			}
		}
		if component, ok := component.(componego.ComponentStop); ok {
			stopComponents = append(stopComponents, component)
		}
	}
	return env, cancelAll
}
