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
	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/ordered-map"
)

func MapTester[T testing.T](
	t testing.TRun[T],
	factory func(cap int) ordered_map.Map[int, float64],
) {
	t.Run("basic features", func(t T) {
		m := factory(0)
		_, ok := m.Get(0)
		require.False(t, ok)
		require.Equal(t, 0, m.Len())
		m.Set(0, 0)
		_, ok = m.Get(0)
		require.True(t, ok)
		m.Set(1, 1)
		m.Set(2, 2)
		require.True(t, m.Has(0))
		require.False(t, m.Has(3))
		require.Equal(t, 3, m.Len())
		require.Equal(t, []int{0, 1, 2}, m.Keys())
		require.Equal(t, []int{2, 1, 0}, m.ReverseKeys())
		require.Equal(t, []float64{0, 1, 2}, m.Values())
		require.Equal(t, []float64{2, 1, 0}, m.ReverseValues())
		m.Set(0, 0.1)
		value, _ := m.Get(0)
		require.Equal(t, 0.1, value)
		m.Remove(0)
		require.False(t, m.Has(0))
		require.Equal(t, 2, m.Len())
		require.Equal(t, []int{1, 2}, m.Keys())
		require.Equal(t, 0, factory(5).Len())
	})
	t.Run("prepend and append", func(t T) {
		m := factory(0)
		m.Set(0, 0)
		m.Prepend(-1, -1)
		m.Prepend(-2, -2)
		m.Append(1, 1)
		m.Append(2, 2)
		require.Equal(t, []float64{-2, -1, 0, 1, 2}, m.Values())
	})
	t.Run("add before and after keys", func(t T) {
		m := factory(0)
		m.AddBefore(1, 1, 1)
		m.AddAfter(7, 7, 7)
		m.AddBefore(0, 0, 1)
		m.AddAfter(2, 2, 1)
		m.AddBefore(3, 3, 7)
		m.AddBefore(5, 5, 7)
		m.AddAfter(4, 4, 3)
		m.AddAfter(6, 6, 5)
		require.Equal(t, []float64{0, 1, 2, 3, 4, 5, 6, 7}, m.Values())
	})
	t.Run("keys position", func(t T) {
		m := factory(0)
		_, ok := m.GetFirstKey()
		require.False(t, ok)
		_, ok = m.GetLastKey()
		require.False(t, ok)
		m.Set(0, 0)
		m.Set(1, 1)
		m.Set(2, 2)
		value, ok := m.GetFirstKey()
		require.True(t, ok)
		require.Equal(t, 0, value)
		value, ok = m.GetLastKey()
		require.True(t, ok)
		require.Equal(t, 2, value)
		_, ok = m.GetPrevKey(0)
		require.False(t, ok)
		value, ok = m.GetPrevKey(1)
		require.True(t, ok)
		require.Equal(t, 0, value)
		_, ok = m.GetNextKey(2)
		require.False(t, ok)
		value, ok = m.GetNextKey(1)
		require.True(t, ok)
		require.Equal(t, 2, value)
	})
	t.Run("swap values", func(t T) {
		m := factory(0)
		m.Set(0, 0)
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)
		require.Equal(t, []int{0, 1, 2, 3}, m.Keys())
		require.Equal(t, []float64{0, 1, 2, 3}, m.Values())
		m.Swap(0, 3)
		m.Swap(2, 1)
		m.Swap(1, 4)
		m.Swap(2, 5)
		require.Equal(t, []int{3, 2, 1, 0}, m.Keys())
		require.Equal(t, []float64{3, 2, 1, 0}, m.Values())
	})
	t.Run("iterate", func(t T) {
		m := factory(0)
		m.Set(1, 1.1)
		m.Set(2, 2.2)
		m.Set(3, 3.3)
		keys := make([]int, 0, m.Len()*2)
		values := make([]float64, 0, m.Len()*2)
		m.Iterate(func(key int, value float64) bool {
			keys = append(keys, key)
			values = append(values, value)
			return true
		})
		m.ReverseIterate(func(key int, value float64) bool {
			keys = append(keys, key)
			values = append(values, value)
			return true
		})
		require.Equal(t, []int{1, 2, 3, 3, 2, 1}, keys)
		require.Equal(t, []float64{1.1, 2.2, 3.3, 3.3, 2.2, 1.1}, values)
		m.Iterate(func(key int, value float64) bool {
			require.Equal(t, 1, key)
			require.Equal(t, 1.1, value)
			return false
		})
		m.ReverseIterate(func(key int, value float64) bool {
			require.Equal(t, 3, key)
			require.Equal(t, 3.3, value)
			return false
		})
	})
	t.Run("convert to simple map", func(t T) {
		m := factory(0)
		m.Set(1, 1.1)
		m.Set(2, 2.2)
		m.Set(3, 3.3)
		require.Equal(t, map[int]float64{
			1: 1.1,
			2: 2.2,
			3: 3.3,
		}, m.ToMap())
		m.Remove(3)
		require.Equal(t, map[int]float64{
			1: 1.1,
			2: 2.2,
		}, m.ToMap())
	})
}
