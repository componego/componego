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
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/internal/utils"
)

func TestKeys(t *testing.T) {
	t.Run("int keys", func(t *testing.T) {
		items := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		expected := []int{1, 2, 3}
		actual := utils.Keys(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("string keys", func(t *testing.T) {
		items := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		expected := []string{"one", "two", "three"}
		actual := utils.Keys(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("empty map", func(t *testing.T) {
		items := map[int]string{}
		var expected []int
		actual := utils.Keys(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("single element map", func(t *testing.T) {
		items := map[string]int{
			"single": 1,
		}
		expected := []string{"single"}
		actual := utils.Keys(items)
		require.ElementsMatch(t, expected, actual)
	})
}

func TestValues(t *testing.T) {
	t.Run("int values", func(t *testing.T) {
		items := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		expected := []int{1, 2, 3}
		actual := utils.Values(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("string values", func(t *testing.T) {
		items := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		expected := []string{"one", "two", "three"}
		actual := utils.Values(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("empty map", func(t *testing.T) {
		items := map[int]string{}
		var expected []string
		actual := utils.Values(items)
		require.ElementsMatch(t, expected, actual)
	})

	t.Run("single element map", func(t *testing.T) {
		items := map[string]int{
			"single": 1,
		}
		expected := []int{1}
		actual := utils.Values(items)
		require.ElementsMatch(t, expected, actual)
	})
}

func TestContains(t *testing.T) {
	t.Run("int slice contains value", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		require.True(t, utils.Contains(items, 3))
	})

	t.Run("int slice does not contain value", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		require.False(t, utils.Contains(items, 6))
	})

	t.Run("string slice contains value", func(t *testing.T) {
		items := []string{"apple", "banana", "cherry"}
		require.True(t, utils.Contains(items, "banana"))
	})

	t.Run("string slice does not contain value", func(t *testing.T) {
		items := []string{"apple", "banana", "cherry"}
		require.False(t, utils.Contains(items, "potato"))
	})

	t.Run("empty slice does not contain value", func(t *testing.T) {
		var items []int
		require.False(t, utils.Contains(items, 1))
	})

	t.Run("single element slice contains value", func(t *testing.T) {
		items := []string{"single"}
		require.True(t, utils.Contains(items, "single"))
	})

	t.Run("single element slice does not contain value", func(t *testing.T) {
		items := []string{"single"}
		require.False(t, utils.Contains(items, "not_single"))
	})
}

func TestReverse(t *testing.T) {
	t.Run("reverse int slice", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		expected := []int{5, 4, 3, 2, 1}
		utils.Reverse(items)
		require.Equal(t, expected, items)
	})

	t.Run("reverse string slice", func(t *testing.T) {
		items := []string{"apple", "banana", "cherry"}
		expected := []string{"cherry", "banana", "apple"}
		utils.Reverse(items)
		require.Equal(t, expected, items)
	})

	t.Run("reverse empty slice", func(t *testing.T) {
		var items []int
		var expected []int
		utils.Reverse(items)
		require.Equal(t, expected, items)
	})

	t.Run("reverse single element slice", func(t *testing.T) {
		items := []string{"single"}
		expected := []string{"single"}
		utils.Reverse(items)
		require.Equal(t, expected, items)
	})

	t.Run("reverse slice with even number of elements", func(t *testing.T) {
		items := []int{1, 2, 3, 4}
		expected := []int{4, 3, 2, 1}
		utils.Reverse(items)
		require.Equal(t, expected, items)
	})
}

func TestCopy(t *testing.T) {
	t.Run("copy int slice", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		expected := []int{1, 2, 3, 4, 5}
		result := utils.Copy(items)
		require.Equal(t, expected, result)
		require.NotSame(t, items, result)
	})

	t.Run("copy string slice", func(t *testing.T) {
		items := []string{"apple", "banana", "cherry"}
		expected := []string{"apple", "banana", "cherry"}
		result := utils.Copy(items)
		require.Equal(t, expected, result)
		require.NotSame(t, items, result)
	})

	t.Run("copy empty slice", func(t *testing.T) {
		items := make([]int, 0)
		expected := make([]int, 0)
		result := utils.Copy(items)
		require.Equal(t, expected, result)
		require.NotSame(t, items, result)
	})

	t.Run("copy single element slice", func(t *testing.T) {
		items := []string{"single"}
		expected := []string{"single"}
		result := utils.Copy(items)
		require.Equal(t, expected, result)
		require.NotSame(t, items, result)
	})

	t.Run("copy slice with even number of elements", func(t *testing.T) {
		items := []int{1, 2, 3, 4}
		expected := []int{1, 2, 3, 4}
		result := utils.Copy(items)
		require.Equal(t, expected, result)
		require.NotSame(t, items, result)
	})
}
