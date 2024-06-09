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
	"strconv"
	"testing"

	"github.com/componego/componego/libs/vendor-proxy"
)

func TestVendorProxy(t *testing.T) {
	VendorProxyTester[*testing.T](t, createInstance)
}

func BenchmarkVendorProxy(b *testing.B) {
	instance := createInstance()
	_ = instance.AddFunction("reflectFunc", func(_ int, _ any, _ ...float64) (bool, error) {
		return true, nil
	})
	_ = instance.AddFunction("contextFunc", func(ctx vendor_proxy.Context, args ...any) (any, error) {
		_ = args[0].(int)
		_ = args[1]
		_ = args[2].([]float64)
		ctx.Validated()
		return true, nil
	})
	arg1 := 1
	arg2 := "text"
	arg3 := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	b.Run("reflect call", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = instance.CallFunction("reflectFunc", arg1, arg2, arg3)
		}
	})
	b.Run("context call", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = instance.CallFunction("contextFunc", arg1, arg2, arg3)
		}
	})
}

func createInstance() vendor_proxy.Proxy {
	index := 0
	for {
		name := "temp-vendor-proxy-" + strconv.Itoa(index)
		if !vendor_proxy.Has(name) {
			return vendor_proxy.Get(name)
		}
		index++
	}
}
