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
	"reflect"
	"strings"
	"testing"

	"github.com/componego/componego/libs/vendor-proxy"
)

func VendorProxyTester(t *testing.T, factory func() vendor_proxy.Proxy) {
	errCustom := errors.New("error text")
	t.Run("using reflect", func(t *testing.T) {
		instance := factory()
		addFunction(t, nil, instance, "func1", func(arg int) int {
			return arg * 2
		})
		equal(t, 3*2, callFunction(t, nil, instance, "func1", 3))
		addFunction(t, vendor_proxy.ErrFunctionExists, instance, "func1", func() {})
		addFunction(t, vendor_proxy.ErrInvalidArgument, instance, "func2", 123)
		addFunction(t, nil, instance, "func2", func() {})
		equal(t, nil, callFunction(t, nil, instance, "func2"))
		addFunction(t, nil, instance, "func3", func(arg1 int, arg2 ...string) (int, string) {
			return arg1, strings.Join(arg2, "-")
		})
		result2 := callFunction(t, nil, instance, "func3", 1, []string{"1", "2", "3"}).([]any)
		equal(t, 1, result2[0])
		equal(t, "1-2-3", result2[1])
		callFunction(t, vendor_proxy.ErrInvalidArgument, instance, "func3", "1")
		callFunction(t, vendor_proxy.ErrWrongCountArguments, instance, "func3", 1, []string{"1", "2", "3"}, nil)
		addFunction(t, nil, instance, "func4", func() (float64, error) {
			return 1.1, errCustom
		})
		equal(t, nil, callFunction(t, errCustom, instance, "func4"))
		addFunction(t, nil, instance, "func5", func() (float64, error) {
			return 1.1, nil
		})
		equal(t, 1.1, callFunction(t, nil, instance, "func5"))
	})
	t.Run("using context", func(t *testing.T) {
		instance := factory()
		addFunction(t, nil, instance, "func1", func(ctx vendor_proxy.Context, args ...any) (any, error) {
			arg := args[0].(int)
			ctx.Validated()
			return arg * 2, nil
		})
		equal(t, 2*2, callFunction(t, nil, instance, "func1", 2))
		callFunction(t, vendor_proxy.ErrInvalidArgument, instance, "func1", "2")
		addFunction(t, nil, instance, "func2", func(_ vendor_proxy.Context, _ ...any) (any, error) {
			return nil, nil
		})
		callFunction(t, vendor_proxy.ErrArgumentsNotValidated, instance, "func2")
		addFunction(t, nil, instance, "func3", func(_ vendor_proxy.Context, args ...any) (any, error) {
			_ = args[0].(float64)
			return nil, nil
		})
		callFunction(t, vendor_proxy.ErrInvalidArgument, instance, "func3", nil)
		addFunction(t, nil, instance, "func4", func(ctx vendor_proxy.Context, _ ...any) (any, error) {
			ctx.Validated()
			ctx.ValidationFailed("error message")
			return nil, nil
		})
		callFunction(t, vendor_proxy.ErrInvalidArgument, instance, "func4")
		addFunction(t, nil, instance, "func5", func(ctx vendor_proxy.Context, _ ...any) (any, error) {
			ctx.Validated()
			return 1.1, errCustom
		})
		equal(t, nil, callFunction(t, errCustom, instance, "func5"))
		addFunction(t, nil, instance, "func6", func(ctx vendor_proxy.Context, _ ...any) (any, error) {
			ctx.Validated()
			return 1.1, nil
		})
		equal(t, 1.1, callFunction(t, nil, instance, "func6"))
	})
}

func equal(t *testing.T, expected any, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		formatFatal(t, expected, actual)
	}
}

func addFunction(t *testing.T, expectedErr error, instance vendor_proxy.Proxy, name string, function any) {
	actualErr := instance.AddFunction(name, function)
	if !errors.Is(actualErr, expectedErr) {
		formatFatal(t, expectedErr, actualErr)
	}
}

func callFunction(t *testing.T, expectedErr error, instance vendor_proxy.Proxy, name string, args ...any) any {
	result, actualErr := instance.CallFunction(name, args...)
	if !errors.Is(actualErr, expectedErr) {
		formatFatal(t, expectedErr, actualErr)
	}
	return result
}

func formatFatal(t *testing.T, expected any, actual any) {
	t.Fatalf("\nexpected: %+v\nactual: %+v", expected, actual)
}
