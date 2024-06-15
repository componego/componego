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

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/xerrors"
)

func TestUnwrapAll(t *testing.T) {
	t.Run("single unwrap error", func(t *testing.T) {
		err1 := errors.New("error")
		err2 := &singleUnwrapError{message: "wrapped error", inner: err1}
		result := xerrors.UnwrapAll(err2)
		require.Len(t, result, 2)
		require.Same(t, err2, result[0])
		require.Same(t, err1, result[1])
	})

	t.Run("multiple unwrap error", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := &multipleUnwrapError{message: "wrapped error", inners: []error{err1, err2}}
		result := xerrors.UnwrapAll(err3)
		require.Len(t, result, 3)
		require.Same(t, err3, result[0])
		require.Same(t, err2, result[1])
		require.Same(t, err1, result[2])
	})

	t.Run("mixed unwrap error", func(t *testing.T) {
		err1 := errors.New("error")
		err2 := &singleUnwrapError{message: "wrapped error 1", inner: err1}
		err3 := &multipleUnwrapError{message: "wrapped error 2", inners: []error{err1, err2}}
		result := xerrors.UnwrapAll(err3)
		require.Len(t, result, 4)
		require.Same(t, err3, result[0])
		require.Same(t, err2, result[1])
		require.Same(t, err1, result[2])
		require.Same(t, err1, result[3])
	})

	t.Run("nil Error", func(t *testing.T) {
		var err error
		result := xerrors.UnwrapAll(err)
		require.Len(t, result, 1)
		require.Nil(t, result[0])
	})
}

type singleUnwrapError struct {
	message string
	inner   error
}

func (e *singleUnwrapError) Error() string {
	return e.message
}

func (e *singleUnwrapError) Unwrap() error {
	return e.inner
}

type multipleUnwrapError struct {
	message string
	inners  []error
}

func (e *multipleUnwrapError) Error() string {
	return e.message
}

func (e *multipleUnwrapError) Unwrap() []error {
	return e.inners
}

var (
	_ error                         = (*singleUnwrapError)(nil)
	_ interface{ Unwrap() error }   = (*singleUnwrapError)(nil)
	_ error                         = (*multipleUnwrapError)(nil)
	_ interface{ Unwrap() []error } = (*multipleUnwrapError)(nil)
)
