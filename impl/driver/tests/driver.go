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
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/logger"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/testing/types"
)

func DriverTester[T testing.TRun[T]](
	t testing.TRun[T],
	factory func() driver.Driver,
) {
	modes := map[string]componego.ApplicationMode{
		"production": componego.ProductionMode,
		"developer":  componego.DeveloperMode,
		"test":       componego.TestMode,
		"another":    123,
	}
	for name, mode := range modes {
		t.Run(fmt.Sprintf("in the %s mode", name), func(t T) {
			if p, ok := any(t).(testing.TParallel); ok {
				p.Parallel()
			}
			driverTester[T](t, factory, mode)
		})
	}
}

func driverTester[T testing.TRun[T]](
	t testing.TRun[T],
	factory func() driver.Driver,
	appMode componego.ApplicationMode,
) {
	customErr := errors.New("custom error")

	t.Run("basic test", func(t T) {
		testCases := [...]struct {
			actionErr  error
			actionCode int
			returnErr  error
			returnCode int
		}{
			{
				actionErr:  nil,
				actionCode: componego.SuccessExitCode,
				returnErr:  nil,
				returnCode: componego.SuccessExitCode,
			},
			{
				actionErr:  nil,
				actionCode: componego.ErrorExitCode,
				returnErr:  nil,
				returnCode: componego.ErrorExitCode,
			},
			{
				actionErr:  customErr,
				actionCode: componego.ErrorExitCode,
				returnErr:  customErr,
				returnCode: componego.ErrorExitCode,
			},
			{
				actionErr:  customErr,
				actionCode: componego.SuccessExitCode,
				returnErr:  customErr,
				returnCode: componego.ErrorExitCode,
			},
			{
				actionErr:  nil,
				actionCode: 123,
				returnErr:  nil,
				returnCode: 123,
			},
			{
				actionErr:  customErr,
				actionCode: 123,
				returnErr:  customErr,
				returnCode: 123,
			},
		}
		for _, testCase := range testCases {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationAction(func(_ componego.Environment, _ any) (int, error) {
				return testCase.actionCode, testCase.actionErr
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, testCase.returnCode, exitCode)
				require.ErrorIs(t, err, testCase.returnErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		}
	})

	t.Run("initialization order", func(t T) {
		buffer := &bytes.Buffer{}
		component1Factory := component.NewFactory("component 1", "0.0.1")
		component2Factory := component.NewFactory("component 2", "0.0.1")
		component2Factory.SetComponentInit(func(_ componego.Environment) error {
			logger.LogData(buffer, "componentInit", "component 2")
			return nil
		})
		component2Factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
			logger.LogData(buffer, "componentStop", "component 2")
			return prevErr
		})
		component3Factory := component.NewFactory("component 3", "0.0.1")
		component3Factory.SetComponentInit(func(_ componego.Environment) error {
			logger.LogData(buffer, "componentInit", "component 3")
			return nil
		})
		component3Factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
			logger.LogData(buffer, "componentStop", "component 3")
			return prevErr
		})
		component1Factory.SetComponentComponents(func() ([]componego.Component, error) {
			return []componego.Component{
				component2Factory.Build(),
			}, nil
		})
		component1Factory.SetComponentInit(func(_ componego.Environment) error {
			logger.LogData(buffer, "componentInit", "component 1")
			return nil
		})
		component1Factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
			logger.LogData(buffer, "componentStop", "component 1")
			return prevErr
		})
		appFactory := application.NewFactory("Application Basic Test")
		appFactory.SetApplicationConfigInit(func(_ componego.ApplicationMode, _ any) (map[string]any, error) {
			logger.LogData(buffer, "applicationConfigInit")
			return nil, nil
		})
		appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
			logger.LogData(buffer, "applicationComponents")
			return []componego.Component{
				component1Factory.Build(),
				component3Factory.Build(),
			}, nil
		})
		appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
			logger.LogData(buffer, "applicationDependencies")
			return nil, nil
		})
		appFactory.SetApplicationAction(func(_ componego.Environment, _ any) (int, error) {
			logger.LogData(buffer, "applicationAction")
			return componego.SuccessExitCode, nil
		})
		require.NotPanics(t, func() {
			exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
			require.Equal(t, componego.SuccessExitCode, exitCode)
			require.NoError(t, err)
		})
		require.Equal(t, logger.ExpectedLogData(
			"applicationConfigInit",
			"applicationComponents",
			"applicationDependencies",
			[]string{"componentInit", "component 2"},
			[]string{"componentInit", "component 1"},
			[]string{"componentInit", "component 3"},
			"applicationAction",
			[]string{"componentStop", "component 3"},
			[]string{"componentStop", "component 1"},
			[]string{"componentStop", "component 2"},
		), buffer.String())
		// In a different order.
		buffer.Reset()
		appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
			logger.LogData(buffer, "applicationComponents")
			return []componego.Component{
				component3Factory.Build(),
				component1Factory.Build(),
			}, nil
		})
		require.NotPanics(t, func() {
			exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
			require.Equal(t, componego.SuccessExitCode, exitCode)
			require.NoError(t, err)
		})
		require.Equal(t, logger.ExpectedLogData(
			"applicationConfigInit",
			"applicationComponents",
			"applicationDependencies",
			[]string{"componentInit", "component 3"},
			[]string{"componentInit", "component 2"},
			[]string{"componentInit", "component 1"},
			"applicationAction",
			[]string{"componentStop", "component 1"},
			[]string{"componentStop", "component 2"},
			[]string{"componentStop", "component 3"},
		), buffer.String())
	})

	t.Run("initialize configuration", func(t T) {
		t.Run("return an error", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationConfigInit(func(_ componego.ApplicationMode, _ any) (map[string]any, error) {
				return nil, customErr
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, customErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})

		t.Run("throw panic as a string", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationConfigInit(func(_ componego.ApplicationMode, _ any) (map[string]any, error) {
				panic("panic occurred")
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorContains(t, err, "panic occurred")
			})
		})

		t.Run("throw panic as a custom type", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationConfigInit(func(_ componego.ApplicationMode, _ any) (map[string]any, error) {
				panic(types.BStruct{})
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})
	})

	t.Run("initialize components", func(t T) {
		t.Run("return an error", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				return nil, customErr
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, customErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})

		t.Run("subcomponent returns an error", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				componentFactory := component.NewFactory("component", "0.0.1")
				componentFactory.SetComponentComponents(func() ([]componego.Component, error) {
					return nil, customErr
				})
				return []componego.Component{
					componentFactory.Build(),
				}, nil
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, customErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})

		t.Run("throw panic as a string", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				panic("panic occurred")
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorContains(t, err, "panic occurred")
			})
		})

		t.Run("throw panic as a custom type", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationComponents(func() ([]componego.Component, error) {
				panic(types.BStruct{})
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})
	})

	t.Run("initialize dependencies", func(t T) {
		t.Run("return an error", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
				return nil, customErr
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, customErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})

		t.Run("dependency constructor returns an error", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
				return []componego.Dependency{
					func() (*types.AStruct, error) {
						return nil, customErr
					},
				}, nil
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, customErr)
				require.NotErrorIs(t, err, driver.ErrPanic)
				require.NotErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})

		t.Run("throw panic as a string", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
				panic("panic occurred")
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorContains(t, err, "panic occurred")
			})
		})

		t.Run("throw panic as a custom type", func(t T) {
			appFactory := application.NewFactory("Application Basic Test")
			appFactory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
				panic(types.BStruct{})
			})
			require.NotPanics(t, func() {
				exitCode, err := factory().RunApplication(context.Background(), appFactory.Build(), appMode)
				require.Equal(t, componego.ErrorExitCode, exitCode)
				require.ErrorIs(t, err, driver.ErrPanic)
				require.ErrorIs(t, err, driver.ErrUnknownPanic)
			})
		})
	})
}
