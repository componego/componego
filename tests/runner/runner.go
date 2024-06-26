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

package runner

import (
	"context"
	"io"
	"runtime"
	"sync"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/internal/testing"
)

type stopAction = func(env componego.Environment, prevErr error) error

func CreateTestEnvironment(t testing.T, app componego.Application, options *driver.Options) (componego.Environment, func()) {
	if options == nil {
		options = &driver.Options{}
	}
	if options.AppIO == nil {
		// Ignore output. We don't need output during the test.
		options.AppIO = application.NewIO(nil, io.Discard, io.Discard)
	}
	cancelableCtx, cancelCtx := context.WithCancel(context.Background())
	t.Cleanup(cancelCtx)
	env, err := driver.New(options).CreateEnvironment(cancelableCtx, app, componego.TestMode)
	if err != nil {
		t.Errorf("error when creating an environment for the application: %s", err)
		t.FailNow()
		return nil, nil
	}
	components := env.Components()
	onStopActions := make([]stopAction, 0, len(components)+1)
	mutexOnce := sync.Once{}
	cancelEnv := func() {
		// Cancellation occurs only once.
		mutexOnce.Do(func() {
			defer func() {
				cancelCtx()
				runtime.Gosched()
				if err = driver.ErrorRecoveryOnStop(recover(), err); err != nil {
					t.Errorf("error when stopping the application: %s", err)
					t.FailNow()
				}
			}()
			for _, action := range onStopActions {
				// noinspection ALL
				defer func(action stopAction) {
					err = driver.ErrorRecoveryOnStop(recover(), err)
					err = action(env, err)
				}(action)
			}
			runtime.Gosched()
		})
	}
	t.Cleanup(cancelEnv)
	for _, component := range components {
		if component, ok := component.(componego.ComponentInit); ok {
			if err = component.ComponentInit(env); err != nil {
				t.Errorf("error in component '%s': %s", component.ComponentIdentifier(), err)
				t.FailNow()
				return nil, nil
			}
		}
		if component, ok := component.(componego.ComponentStop); ok {
			onStopActions = append(onStopActions, component.ComponentStop)
		}
	}
	return env, cancelEnv
}
