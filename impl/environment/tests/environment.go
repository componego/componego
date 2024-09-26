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
	"context"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/environment"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/impl/environment/managers/config"
	"github.com/componego/componego/impl/environment/managers/dependency"
	"github.com/componego/componego/internal/system"
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/require"
)

type testCtxKey struct{}

func EnvironmentTester[T testing.TRun[T]](
	t testing.TRun[T],
	factory func(
		context context.Context,
		application componego.Application,
		applicationIO componego.ApplicationIO,
		applicationMode componego.ApplicationMode,
		configProvider componego.ConfigProvider,
		componentProvider componego.ComponentProvider,
		dependencyInvoker componego.DependencyInvoker,
	) componego.Environment,
) {
	components := []componego.Component{
		component.NewFactory("tests:component1", "0.0.1").Build(),
		component.NewFactory("tests:component2", "0.0.2").Build(),
	}
	app := application.NewFactory("Test Application").Build()
	appIO := application.NewIO(system.Stdin, system.Stdout, system.Stderr)
	appMode := componego.ProductionMode
	configProvider, _ := config.NewManager()
	componentProvider, componentsInitializer := component.NewManager()
	require.NoError(t, componentsInitializer(components))
	dependencyInvoker, _ := dependency.NewManager()

	ctxKey := testCtxKey{}

	t.Run("get objects", func(t T) {
		ctx := context.WithValue(context.Background(), ctxKey, 123)
		env := factory(ctx, app, appIO, appMode, configProvider, componentProvider, dependencyInvoker)
		require.Same(t, app, env.Application())
		require.Same(t, appIO, env.ApplicationIO())
		require.Same(t, configProvider, env.ConfigProvider())
		require.Same(t, dependencyInvoker, env.DependencyInvoker())
		require.Equal(t, 123, env.GetContext().Value(ctxKey))
		require.Equal(t, appMode, env.ApplicationMode())
		require.Equal(t, componentProvider.Components(), env.Components())
	})

	t.Run("set context", func(t T) {
		env := factory(context.Background(), app, appIO, appMode, configProvider, componentProvider, dependencyInvoker)
		err := env.SetContext(context.Background())
		require.ErrorIs(t, err, environment.ErrInvalidParentContext)
		err = env.SetContext(env.GetContext())
		require.NoError(t, err)
		err = env.SetContext(context.WithValue(env.GetContext(), ctxKey, 321))
		require.NoError(t, err)
	})

	t.Run("get environment from context", func(t T) {
		expectedEnv := factory(context.Background(), app, appIO, appMode, configProvider, componentProvider, dependencyInvoker)

		t.Run("without panic", func(t T) {
			actualEnv, err := environment.GetEnvironment(expectedEnv.GetContext())
			require.NoError(t, err)
			require.Same(t, expectedEnv, actualEnv)
			actualEnv, err = environment.GetEnvironment(context.Background())
			require.ErrorIs(t, err, environment.ErrNoEnvironmentInContext)
			require.Nil(t, actualEnv)
		})

		t.Run("with panic", func(t T) {
			require.NotPanics(t, func() {
				actualEnv := environment.GetEnvironmentOrPanic(expectedEnv.GetContext())
				require.Same(t, expectedEnv, actualEnv)
			})
			require.PanicsWithError(t, environment.ErrNoEnvironmentInContext.Error(), func() {
				environment.GetEnvironmentOrPanic(context.Background())
			})
		})
	})
}
