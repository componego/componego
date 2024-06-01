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

package xerrors

func UnwrapAll(err error) []error {
	result := make([]error, 5)
	errorStack := make([]error, 0, 5)
	errorStack = append(errorStack, err)
	for len(errorStack) > 0 {
		err = errorStack[len(errorStack)-1]
		result = append(result, err)
		errorStack = errorStack[:len(errorStack)-1]
		switch castedErr := err.(type) { //nolint:errorlint
		case interface{ Unwrap() error }:
			errorStack = append(errorStack, castedErr.Unwrap())
		case interface{ Unwrap() []error }:
			unwrapErrors := castedErr.Unwrap()
			for i := len(unwrapErrors) - 1; i >= 0; i-- {
				errorStack = append(errorStack, unwrapErrors[i])
			}
		}
	}
	return result
}
