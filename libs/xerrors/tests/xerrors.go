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
	"strings"

	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/libs/xerrors"
)

func XErrorsTester[T testing.T](
	t testing.TRun[T],
	factory func(message string, code string, options ...xerrors.Option) xerrors.XError,
) {
	err1 := errors.New("error1")
	err1x := factory("error1", "E1X")
	err2 := errors.New("error2")
	err2x := factory("error2", "E2X")
	err3 := errors.New("error3")
	err3x := factory("error3", "E3X")
	err1S3 := errors.Join(err1, err3)
	err1xS3x := err1x.WithError(err3x, "E1XS3X")
	err1S2 := errors.Join(err1, err2)
	err1xS2 := err1x.WithError(err2, "E1XS2")
	err1S2S3 := errors.Join(err1S2, err3)
	err1xS2S3x := err1xS2.WithError(err3x, "E1XS2S3X")
	t.Run("comparison with the original error package", func(t T) {
		testCases := [...]struct {
			a1 error
			b1 error
			a2 error
			b2 error
		}{
			{err1, err1, err1x, err1x},          // 1
			{err1, err2, err1x, err2x},          // 2
			{err2, err1, err2x, err1x},          // 3
			{err1, err3, err1x, err3x},          // 4
			{err1, err1S3, err1x, err1xS3x},     // 5
			{err1S3, err1, err1xS3x, err1x},     // 6
			{err1S3, err2, err1xS3x, err2x},     // 7
			{err1S2, err2, err1xS2, err2},       // 8
			{err2, err1S2, err2, err1xS2},       // 9
			{err1S2S3, err3, err1xS2S3x, err3x}, // 10
			{err1S2, err3, err1xS2, err3x},      // 11
		}
		for i, testCase := range testCases {
			if errors.Is(testCase.a1, testCase.b1) != errors.Is(testCase.a2, testCase.b2) {
				t.Errorf("#%d failed", i+1)
				t.FailNow()
			}
		}
	})
	t.Run("comparison with the expected result", func(t T) {
		testCases := [...]struct {
			a      error
			b      error
			result bool
		}{
			{err1x, err1x, true},       // 1
			{err1x, err2x, false},      // 2
			{err1x, err1, false},       // 3
			{err1, err1x, false},       // 4
			{err1xS3x, err1x, true},    // 5
			{err1xS3x, err3x, true},    // 6
			{err1xS3x, err2x, false},   // 7
			{err1xS2, err2, true},      // 8
			{err2, err1xS2, false},     // 9
			{err1xS2, err2, true},      // 10
			{err1xS2S3x, err1x, true},  // 11
			{err1xS2S3x, err2, true},   // 12
			{err1xS2S3x, err3x, true},  // 13
			{err1xS2S3x, err3, false},  // 14
			{err3x, err1xS2S3x, false}, // 15
		}
		for i, testCase := range testCases {
			if errors.Is(testCase.a, testCase.b) != testCase.result {
				t.Errorf("#%d failed", i+1)
				t.FailNow()
			}
		}
	})
	t.Run("error text comparison", func(t T) {
		testCases2 := [...]struct {
			a error
			b error
		}{
			{err1, err1x},          // 1
			{err1S3, err1xS3x},     // 2
			{err1S2, err1xS2},      // 3
			{err1S2S3, err1xS2S3x}, // 4
		}
		for i, testCase := range testCases2 {
			if testCase.a.Error() != convertErrorMessage(testCase.b.Error()) {
				t.Errorf("#%d failed", i+1)
				t.FailNow()
			}
		}
	})
}

func GetOptionsByKey(t testing.T, err error, key string) any {
	for _, err = range xerrors.UnwrapAll(err) {
		// noinspection ALL
		xError, ok := err.(xerrors.XError) //nolint:errorlint
		if !ok {
			continue
		}
		for _, option := range xError.ErrorOptions() {
			if option.Key() == key {
				return option.Value()
			}
		}
	}
	t.Errorf("failed to get error option by key '%s'", key)
	t.FailNow()
	return nil
}

func convertErrorMessage(errMessage string) string {
	lines := strings.Split(errMessage, "\n")
	for i, line := range lines {
		index := strings.Index(line, " (")
		if index != -1 {
			lines[i] = line[:index]
		}
	}
	return strings.Join(lines, "\n")
}
