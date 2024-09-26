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
	"testing"

	"github.com/componego/componego/libs/ordered-map"
)

func TestOrderedMap(t *testing.T) {
	MapTester[*testing.T](t, ordered_map.New[int, float64])
}

func BenchmarkOrderedMap(b *testing.B) {
	b.Run("set", func(b *testing.B) {
		m := ordered_map.New[int, int](b.N)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.Set(n, n)
		}
	})
	b.Run("get", func(b *testing.B) {
		m := ordered_map.New[int, int](1)
		m.Set(0, 0)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.Get(0)
		}
	})
	b.Run("has", func(b *testing.B) {
		m := ordered_map.New[int, int](1)
		m.Set(0, 0)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.Has(0)
		}
	})
	b.Run("prepend", func(b *testing.B) {
		m := ordered_map.New[int, int](b.N)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.Prepend(n, n)
		}
	})
	b.Run("append", func(b *testing.B) {
		m := ordered_map.New[int, int](b.N)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.Append(n, n)
		}
	})
	b.Run("add before", func(b *testing.B) {
		m := ordered_map.New[int, int](b.N)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.AddBefore(n, n, n-1)
		}
	})
	b.Run("add after", func(b *testing.B) {
		m := ordered_map.New[int, int](b.N)
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			m.AddAfter(n, n, n-1)
		}
	})
}
