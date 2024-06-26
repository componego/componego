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

package tests

import (
	"errors"
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/internal/testing/require"
)

func TestExitWrapper(t *testing.T) {
	errCustom := errors.New("custom error")
	exitCode, err := application.ExitWrapper(errCustom)
	require.Equal(t, componego.ErrorExitCode, exitCode)
	require.ErrorIs(t, err, errCustom)
	exitCode, err = application.ExitWrapper(nil)
	require.Equal(t, componego.SuccessExitCode, exitCode)
	require.NoError(t, err)
}
