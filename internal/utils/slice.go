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

package utils

func Keys[K comparable, V any](items map[K]V) []K {
	result := make([]K, 0, len(items))
	for key := range items {
		result = append(result, key)
	}
	return result
}

func Values[K comparable, V any](items map[K]V) []V {
	result := make([]V, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	return result
}

func Contains[T comparable](items []T, value T) bool {
	for _, item := range items {
		if value == item {
			return true
		}
	}
	return false
}

func Reverse[T any](items []T) {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}

func Copy[T any](items []T) []T {
	result := make([]T, len(items))
	copy(result, items)
	return result
}
