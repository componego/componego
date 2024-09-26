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
	"errors"
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment"
	"github.com/componego/componego/impl/environment/managers/component"
	"github.com/componego/componego/internal/testing/require"
)

func TestFactory(t *testing.T) {
	t.Run("check build method", func(t *testing.T) {
		factory := component.NewFactory("test", "0.0.1")
		require.NotSame(t, factory.Build(), factory.Build())
	})

	t.Run("compare identifier and version", func(t *testing.T) {
		factory := component.NewFactory("test", "0.0.1")
		componentItem1 := factory.Build()
		require.Equal(t, "test", componentItem1.ComponentIdentifier())
		require.Equal(t, "0.0.1", componentItem1.ComponentVersion())
		factory.SetComponentIdentifier("test2")
		factory.SetComponentVersion("0.0.2")
		require.Equal(t, "test", componentItem1.ComponentIdentifier())
		require.Equal(t, "0.0.1", componentItem1.ComponentVersion())
		componentItem2 := factory.Build()
		require.Equal(t, "test2", componentItem2.ComponentIdentifier())
		require.Equal(t, "0.0.2", componentItem2.ComponentVersion())
	})

	t.Run("compare components", func(t *testing.T) {
		factory := component.NewFactory("test", "0.0.1")
		componentItem1 := factory.Build()
		require.NotPanics(t, func() {
			actualComponents, actualErr := componentItem1.(componego.ComponentComponents).ComponentComponents()
			require.Len(t, actualComponents, 0)
			require.NoError(t, actualErr)
		})
		expectedComponents := []componego.Component{nil, nil}
		expectedErrorForComponents := errors.New("error 1")
		factory.SetComponentComponents(func() ([]componego.Component, error) {
			return expectedComponents, expectedErrorForComponents
		})
		componentItem2 := factory.Build()
		require.NotPanics(t, func() {
			actualDependencies, actualErr := componentItem2.(componego.ComponentComponents).ComponentComponents()
			require.Equal(t, expectedComponents, actualDependencies)
			require.ErrorIs(t, actualErr, expectedErrorForComponents)
		})
	})

	t.Run("compare dependencies", func(t *testing.T) {
		factory := component.NewFactory("test", "0.0.1")
		componentItem1 := factory.Build()
		require.NotPanics(t, func() {
			actualDependencies, actualErr := componentItem1.(componego.ComponentDependencies).ComponentDependencies()
			require.Len(t, actualDependencies, 0)
			require.NoError(t, actualErr)
		})
		expectedDependencies := []componego.Dependency{nil, nil}
		expectedErrorForDependencies := errors.New("error 1")
		factory.SetComponentDependencies(func() ([]componego.Dependency, error) {
			return expectedDependencies, expectedErrorForDependencies
		})
		componentItem2 := factory.Build()
		require.NotPanics(t, func() {
			actualDependencies, actualErr := componentItem2.(componego.ComponentDependencies).ComponentDependencies()
			require.Equal(t, expectedDependencies, actualDependencies)
			require.ErrorIs(t, actualErr, expectedErrorForDependencies)
		})
	})

	t.Run("run init and stop methods", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := errors.New("error 3")
		expectedEnv := environment.New(context.Background(), nil, nil, componego.ProductionMode, nil, nil, nil)
		factory := component.NewFactory("test", "0.0.1")
		componentItem1 := factory.Build()
		require.NotPanics(t, func() {
			require.NoError(t, componentItem1.(componego.ComponentInit).ComponentInit(expectedEnv))
			require.NoError(t, componentItem1.(componego.ComponentStop).ComponentStop(expectedEnv, nil))
			require.ErrorIs(t, componentItem1.(componego.ComponentStop).ComponentStop(expectedEnv, err2), err2)
		})
		factory.SetComponentInit(func(actualEnv componego.Environment) error {
			require.Same(t, expectedEnv, actualEnv)
			return err1
		})
		factory.SetComponentStop(func(actualEnv componego.Environment, prevErr error) error {
			require.Same(t, expectedEnv, actualEnv)
			require.ErrorIs(t, prevErr, err2)
			return err3
		})
		componentItem2 := factory.Build()
		require.NotPanics(t, func() {
			require.ErrorIs(t, componentItem2.(componego.ComponentInit).ComponentInit(expectedEnv), err1)
			require.ErrorIs(t, componentItem2.(componego.ComponentStop).ComponentStop(expectedEnv, err2), err3)
			require.NotErrorIs(t, componentItem2.(componego.ComponentStop).ComponentStop(expectedEnv, err2), err2)
		})
		factory.SetComponentStop(func(_ componego.Environment, prevErr error) error {
			return errors.Join(prevErr, err3)
		})
		componentItem3 := factory.Build()
		require.NotPanics(t, func() {
			require.ErrorIs(t, componentItem3.(componego.ComponentStop).ComponentStop(expectedEnv, err2), err3)
			require.ErrorIs(t, componentItem3.(componego.ComponentStop).ComponentStop(expectedEnv, err2), err2)
		})
	})
}
