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

package tests

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing/logger"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/tests/runner"
)

func TestCreateTestEnvironment(t *testing.T) {
	error1 := errors.New("error 1")
	error2 := errors.New("error 2")

	t.Run("create environment", func(t *testing.T) {
		t.Run("without errors", func(t *testing.T) {
			app := application.NewFactory("test").Build()
			env, cancelEnv := runner.CreateTestEnvironment(t, app, nil)
			require.Same(t, app, env.Application())
			require.NotNil(t, cancelEnv)
			require.Equal(t, componego.TestMode, env.ApplicationMode())
		})

		t.Run("with errors", func(t *testing.T) {
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
				panic("any panic")
				return nil, nil //nolint:govet
			})
			app := appFactory.Build()
			mockedT := &testingMock{}
			env, cancelEnv := runner.CreateTestEnvironment(mockedT, app, nil)
			require.Nil(t, env)
			require.Nil(t, cancelEnv)
			require.True(t, mockedT.IsFailed)
		})
	})

	t.Run("cancel environment", func(t *testing.T) {
		t.Run("it must be cancelled in any case", func(t *testing.T) {
			canceled := false
			appFactory := application.NewFactory("test")
			options := &runner.TestOptions{
				OnApplicationStop: func(_ componego.Application, returnErr error) error {
					require.NoError(t, returnErr)
					canceled = true
					return returnErr
				},
			}
			t.Run("internal test 1", func(t *testing.T) {
				require.False(t, canceled)
				_, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), options)
				require.False(t, canceled)
				cancelEnv()
				require.True(t, canceled)
			})
			require.True(t, canceled)
			canceled = false
			t.Run("internal test 2", func(t *testing.T) {
				require.False(t, canceled)
				_, _ = runner.CreateTestEnvironment(t, appFactory.Build(), options)
				require.False(t, canceled)
			})
			require.True(t, canceled)
		})

		t.Run("it must be cancelled once", func(t *testing.T) {
			canceled := false
			app := application.NewFactory("test").Build()
			_, cancelEnv := runner.CreateTestEnvironment(t, app, &runner.TestOptions{
				OnApplicationStop: func(_ componego.Application, returnErr error) error {
					require.NoError(t, returnErr)
					canceled = true
					return returnErr
				},
			})
			require.False(t, canceled)
			cancelEnv()
			require.True(t, canceled)
			canceled = false
			cancelEnv()
			require.False(t, canceled)
		})
	})

	t.Run("custom driver and application run mode passed to the environment", func(t *testing.T) {
		created := false
		appFactory := application.NewFactory("test")
		app := appFactory.Build()
		appMode := componego.ProductionMode
		customDriver := driver.New(nil)
		env, _ := runner.CreateTestEnvironment(t, app, &runner.TestOptions{
			Driver: customDriver,
			EnvironmentFactory: func(internalApp componego.Application, internalDriver driver.Driver) (componego.Environment, func() error, error) {
				require.Same(t, app, internalApp)
				require.Same(t, customDriver, internalDriver)
				env, cancelEnv, err := internalDriver.CreateEnvironment(context.Background(), internalApp, appMode)
				created = true
				return env, cancelEnv, err
			},
		})
		require.True(t, created)
		require.Same(t, app, env.Application())
		require.Equal(t, appMode, env.ApplicationMode())
	})

	t.Run("application with components", func(t *testing.T) {
		t.Run("initialize and stop components without errors", func(t *testing.T) {
			buffer := &bytes.Buffer{}
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				component1 := component.NewFactory("component1", "0.0.1")
				component1.SetComponentInit(func(_ componego.Environment) error {
					logger.LogData(buffer, "component1 init")
					return nil
				})
				component1.SetComponentStop(func(_ componego.Environment, prevErr error) error {
					logger.LogData(buffer, "component1 stop")
					return prevErr
				})
				component2 := component.NewFactory("component2", "0.0.1")
				component2.SetComponentInit(func(_ componego.Environment) error {
					logger.LogData(buffer, "component2 init")
					return nil
				})
				component2.SetComponentStop(func(_ componego.Environment, prevErr error) error {
					logger.LogData(buffer, "component2 stop")
					return prevErr
				})
				return []componego.Component{
					component1.Build(),
					component2.Build(),
				}, nil
			})
			env, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), nil)
			require.NotNil(t, env)
			require.Equal(t,
				logger.ExpectedLogData("component1 init", "component2 init"),
				buffer.String(),
			)
			buffer.Reset()
			cancelEnv()
			require.Equal(t,
				logger.ExpectedLogData("component2 stop", "component1 stop"),
				buffer.String(),
			)
		})

		t.Run("initialize and stop components with errors", func(t *testing.T) {
			buffer := &bytes.Buffer{}
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				component1 := component.NewFactory("component1", "0.0.1")
				component1.SetComponentInit(func(_ componego.Environment) error {
					logger.LogData(buffer, "component1 init")
					return nil
				})
				component1.SetComponentStop(func(_ componego.Environment, prevErr error) error {
					logger.LogData(buffer, "component1 stop")
					require.ErrorIs(t, prevErr, error1)
					return errors.Join(prevErr, error2)
				})
				component2 := component.NewFactory("component2", "0.0.1")
				component2.SetComponentInit(func(_ componego.Environment) error {
					return error1
				})
				return []componego.Component{
					component1.Build(),
					component2.Build(),
				}, nil
			})
			env, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), &runner.TestOptions{
				OnApplicationStop: func(_ componego.Application, returnErr error) error {
					require.ErrorIs(t, returnErr, error1)
					require.ErrorIs(t, returnErr, error2)
					logger.LogData(buffer, "application stop")
					return nil
				},
			})
			require.Nil(t, env)
			require.Nil(t, cancelEnv)
			require.Equal(t,
				logger.ExpectedLogData("component1 init", "component1 stop", "application stop"),
				buffer.String(),
			)
		})

		t.Run("initialize and stop components with panic", func(t *testing.T) {
			buffer := &bytes.Buffer{}
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				component1 := component.NewFactory("component1", "0.0.2")
				component1.SetComponentInit(func(_ componego.Environment) error {
					logger.LogData(buffer, "component1 init")
					return nil
				})
				component1.SetComponentStop(func(_ componego.Environment, prevErr error) error {
					logger.LogData(buffer, "component1 stop")
					require.ErrorIs(t, prevErr, error1)
					panic(errors.Join(prevErr, error2))
					return nil //nolint:govet
				})
				component2 := component.NewFactory("component2", "0.0.2")
				component2.SetComponentInit(func(_ componego.Environment) error {
					panic(error1)
					return nil //nolint:govet
				})
				return []componego.Component{
					component1.Build(),
					component2.Build(),
				}, nil
			})
			env, cancelEnv := runner.CreateTestEnvironment(t, appFactory.Build(), &runner.TestOptions{
				OnApplicationStop: func(_ componego.Application, returnErr error) error {
					require.ErrorIs(t, returnErr, error1)
					require.ErrorIs(t, returnErr, error2)
					logger.LogData(buffer, "application stop")
					return nil
				},
			})
			require.Nil(t, env)
			require.Nil(t, cancelEnv)
			require.Equal(t,
				logger.ExpectedLogData("component1 init", "component1 stop", "application stop"),
				buffer.String(),
			)
		})
	})

	t.Run("test hooks", func(t *testing.T) {
		t.Run("on application stop", func(t *testing.T) {
			t.Run("the application stops without errors", func(t *testing.T) {
				hookTriggered := false
				appFactory := application.NewFactory("test")
				app := appFactory.Build()
				mockedT := &testingMock{}
				env, cancelEnv := runner.CreateTestEnvironment(mockedT, app, &runner.TestOptions{
					OnApplicationStop: func(internalApp componego.Application, returnErr error) error {
						require.Same(t, app, internalApp)
						require.NoError(t, returnErr)
						require.False(t, hookTriggered)
						hookTriggered = true
						return returnErr
					},
				})
				require.Same(t, app, env.Application())
				require.False(t, hookTriggered)
				cancelEnv()
				require.True(t, hookTriggered)
				require.False(t, mockedT.IsFailed)
			})

			t.Run("application stops with an error and fails tests", func(t *testing.T) {
				hookTriggered := false
				appFactory := application.NewFactory("test")
				appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
					componentFactory := component.NewFactory("component1", "0.0.1")
					componentFactory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
						require.NoError(t, prevErr)
						return error1
					})
					return []componego.Component{
						componentFactory.Build(),
					}, nil
				})
				app := appFactory.Build()
				mockedT := &testingMock{}
				env, cancelEnv := runner.CreateTestEnvironment(mockedT, app, &runner.TestOptions{
					OnApplicationStop: func(internalApp componego.Application, returnErr error) error {
						require.Same(t, app, internalApp)
						require.ErrorIs(t, returnErr, error1)
						require.False(t, hookTriggered)
						hookTriggered = true
						return error1
					},
				})
				require.Same(t, app, env.Application())
				require.False(t, hookTriggered)
				cancelEnv()
				require.True(t, hookTriggered)
				require.True(t, mockedT.IsFailed)
			})

			t.Run("application stops with an error and does not fail tests", func(t *testing.T) {
				hookTriggered := false
				appFactory := application.NewFactory("test")
				appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
					componentFactory := component.NewFactory("component1", "0.0.1")
					componentFactory.SetComponentInit(func(_ componego.Environment) error {
						return error1
					})
					return []componego.Component{
						componentFactory.Build(),
					}, nil
				})
				app := appFactory.Build()
				mockedT := &testingMock{}
				runner.CreateTestEnvironment(mockedT, app, &runner.TestOptions{
					OnApplicationStop: func(internalApp componego.Application, returnErr error) error {
						require.Same(t, app, internalApp)
						require.ErrorIs(t, returnErr, error1)
						require.False(t, hookTriggered)
						hookTriggered = true
						return nil // This test will not fail.
					},
				})
				require.True(t, hookTriggered)
				require.False(t, mockedT.IsFailed)
			})
		})

		t.Run("on component init", func(t *testing.T) {
			component1 := component.NewFactory("component1", "0.0.1").Build()
			component2Factory := component.NewFactory("component2", "0.0.1")
			component2Factory.SetComponentInit(func(_ componego.Environment) error {
				return error2
			})
			component2 := component2Factory.Build()
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				return []componego.Component{
					component1,
					component2,
				}, nil
			})
			app := appFactory.Build()
			component1HookTriggered := false
			component2HookTriggered := false
			applicationStopHookTriggered := false
			env, cancelEnv := runner.CreateTestEnvironment(t, app, &runner.TestOptions{
				OnComponentInit: func(env componego.Environment, component componego.ComponentInit, returnErr error) {
					require.Same(t, app, env.Application())
					if component1HookTriggered {
						require.Same(t, component2, component)
						require.ErrorIs(t, returnErr, error2)
						require.False(t, component2HookTriggered)
						component2HookTriggered = true
						return
					}
					require.Same(t, component1, component)
					require.NoError(t, returnErr)
					require.False(t, component1HookTriggered)
					component1HookTriggered = true
				},
				OnApplicationStop: func(_ componego.Application, _ error) error {
					applicationStopHookTriggered = true
					return nil // This test will not fail.
				},
			})
			require.True(t, component1HookTriggered)
			require.True(t, component2HookTriggered)
			require.True(t, applicationStopHookTriggered) // There was an error during component initialization so the environment was stopped.
			require.Nil(t, env)
			require.Nil(t, cancelEnv)
		})

		t.Run("on component stop", func(t *testing.T) {
			component1Factory := component.NewFactory("component1", "0.0.1")
			component1Factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
				require.ErrorIs(t, prevErr, error2)
				return error1 // instead of error2
			})
			component1 := component1Factory.Build()
			component2Factory := component.NewFactory("component2", "0.0.1")
			component2Factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
				require.NoError(t, prevErr)
				return error2
			})
			component2 := component2Factory.Build()
			appFactory := application.NewFactory("test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				return []componego.Component{
					component1,
					component2,
				}, nil
			})
			app := appFactory.Build()
			component1HookTriggered := false
			component2HookTriggered := false
			applicationStopHookTriggered := false
			_, cancelEnv := runner.CreateTestEnvironment(t, app, &runner.TestOptions{
				OnComponentStop: func(env componego.Environment, component componego.ComponentStop, returnErr error, previousErr error) {
					require.Same(t, app, env.Application())
					if component2HookTriggered {
						require.Same(t, component1, component)
						require.ErrorIs(t, returnErr, error1)
						require.NotErrorIs(t, returnErr, error2) // because the error was replaced instead of merged in the first component.
						require.ErrorIs(t, previousErr, error2)
						require.False(t, component1HookTriggered)
						component1HookTriggered = true
						return
					}
					require.Same(t, component2, component)
					require.ErrorIs(t, returnErr, error2)
					require.NoError(t, previousErr)
					require.False(t, component2HookTriggered)
					component2HookTriggered = true
				},
				OnApplicationStop: func(_ componego.Application, returnErr error) error {
					require.ErrorIs(t, returnErr, error1)
					require.NotErrorIs(t, returnErr, error2) // because the error was replaced instead of merged in the first component.
					applicationStopHookTriggered = true
					return nil
				},
			})
			require.False(t, component1HookTriggered)
			require.False(t, component2HookTriggered)
			require.False(t, applicationStopHookTriggered)
			cancelEnv()
			require.True(t, component1HookTriggered)
			require.True(t, component2HookTriggered)
			require.True(t, applicationStopHookTriggered)
		})
	})
}

type testingMock struct {
	IsFailed bool
}

func (t *testingMock) Cleanup(_ func()) {
	// empty method.
}

func (t *testingMock) Errorf(_ string, _ ...any) {
	// empty method.
}

func (t *testingMock) FailNow() {
	t.IsFailed = true
}
