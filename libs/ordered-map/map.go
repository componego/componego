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

package ordered_map

type Map[K comparable, V any] interface {
	Set(key K, value V)
	Has(key K) bool
	Get(key K) (V, bool)
	GetFirstKey() (K, bool)
	GetNextKey(key K) (K, bool)
	GetLastKey() (K, bool)
	GetPrevKey(key K) (K, bool)
	Remove(key K)
	Prepend(key K, value V)
	Append(key K, value V)
	AddBefore(key K, value V, beforeKey K)
	AddAfter(key K, value V, afterKey K)
	Keys() []K
	ReverseKeys() []K
	Values() []V
	ReverseValues() []V
	Iterate(fn func(key K, value V) bool)
	ReverseIterate(fn func(key K, value V) bool)
	Swap(key1 K, key2 K)
	Len() int
	ToMap() map[K]V
}

type item[K comparable, V any] struct {
	key        K
	value      V
	prev, next *item[K, V]
}

type _map[K comparable, V any] struct {
	items               map[K]*item[K, V]
	firstItem, lastItem *item[K, V]
}

func New[K comparable, V any](cap int) Map[K, V] {
	return &_map[K, V]{
		items: make(map[K]*item[K, V], cap),
	}
}

func (m *_map[K, V]) Set(key K, value V) {
	if oldItem, ok := m.items[key]; ok {
		oldItem.value = value
		return
	}
	newItem := &item[K, V]{
		key:   key,
		value: value,
	}
	if len(m.items) == 0 {
		m.firstItem = newItem
	} else {
		newItem.prev = m.lastItem
		m.lastItem.next = newItem
	}
	m.lastItem = newItem
	m.items[key] = newItem
}

func (m *_map[K, V]) Has(key K) bool {
	_, ok := m.items[key]
	return ok
}

func (m *_map[K, V]) Get(key K) (V, bool) {
	if resultItem, ok := m.items[key]; ok {
		return resultItem.value, true
	}
	return *new(V), false
}

func (m *_map[K, V]) GetFirstKey() (K, bool) {
	if m.firstItem != nil {
		return m.firstItem.key, true
	}
	return *new(K), false
}

func (m *_map[K, V]) GetNextKey(key K) (K, bool) {
	if resultItem, ok := m.items[key]; ok && resultItem.next != nil {
		return resultItem.next.key, true
	}
	return *new(K), false
}

func (m *_map[K, V]) GetLastKey() (K, bool) {
	if m.lastItem != nil {
		return m.lastItem.key, true
	}
	return *new(K), false
}

func (m *_map[K, V]) GetPrevKey(key K) (K, bool) {
	if resultItem, ok := m.items[key]; ok && resultItem.prev != nil {
		return resultItem.prev.key, true
	}
	return *new(K), false
}

func (m *_map[K, V]) Remove(key K) {
	oldItem, ok := m.items[key]
	if !ok {
		return
	}
	if oldItem.prev != nil {
		oldItem.prev.next = oldItem.next
	}
	if oldItem.next != nil {
		oldItem.next.prev = oldItem.prev
	}
	if m.firstItem == oldItem {
		m.firstItem = oldItem.next
	}
	if m.lastItem == oldItem {
		m.lastItem = oldItem.prev
	}
	delete(m.items, key)
}

func (m *_map[K, V]) Prepend(key K, value V) {
	if len(m.items) == 0 {
		m.Set(key, value)
	} else {
		m.addBeforeInternal(key, value, m.firstItem)
	}
}

func (m *_map[K, V]) Append(key K, value V) {
	if len(m.items) == 0 {
		m.Set(key, value)
	} else {
		m.addAfterInternal(key, value, m.lastItem)
	}
}

func (m *_map[K, V]) AddBefore(key K, value V, beforeKey K) {
	if key == beforeKey {
		m.Set(key, value)
		return
	}
	if beforeItem, ok := m.items[beforeKey]; ok {
		m.addBeforeInternal(key, value, beforeItem)
	} else {
		m.Prepend(key, value)
	}
}

func (m *_map[K, V]) AddAfter(key K, value V, afterKey K) {
	if key == afterKey {
		m.Set(key, value)
		return
	}
	if afterItem, ok := m.items[afterKey]; ok {
		m.addAfterInternal(key, value, afterItem)
	} else {
		m.Append(key, value)
	}
}

func (m *_map[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.items))
	currentItem := m.firstItem
	for currentItem != nil {
		keys = append(keys, currentItem.key)
		currentItem = currentItem.next
	}
	return keys
}

func (m *_map[K, V]) ReverseKeys() []K {
	keys := make([]K, 0, len(m.items))
	currentItem := m.lastItem
	for currentItem != nil {
		keys = append(keys, currentItem.key)
		currentItem = currentItem.prev
	}
	return keys
}

func (m *_map[K, V]) Values() []V {
	values := make([]V, 0, len(m.items))
	currentItem := m.firstItem
	for currentItem != nil {
		values = append(values, currentItem.value)
		currentItem = currentItem.next
	}
	return values
}

func (m *_map[K, V]) ReverseValues() []V {
	values := make([]V, 0, len(m.items))
	currentItem := m.lastItem
	for currentItem != nil {
		values = append(values, currentItem.value)
		currentItem = currentItem.prev
	}
	return values
}

func (m *_map[K, V]) Iterate(fn func(key K, value V) bool) {
	currentItem := m.firstItem
	for currentItem != nil {
		if !fn(currentItem.key, currentItem.value) {
			return
		}
		currentItem = currentItem.next
	}
}

func (m *_map[K, V]) ReverseIterate(fn func(key K, value V) bool) {
	currentItem := m.lastItem
	for currentItem != nil {
		if !fn(currentItem.key, currentItem.value) {
			return
		}
		currentItem = currentItem.prev
	}
}

func (m *_map[K, V]) Swap(key1 K, key2 K) {
	item1, ok1 := m.items[key1]
	item2, ok2 := m.items[key2]
	if ok1 && ok2 {
		item1.value, item2.value = item2.value, item1.value
		item1.key, item2.key = item2.key, item1.key
	}
}

func (m *_map[K, V]) Len() int {
	return len(m.items)
}

func (m *_map[K, V]) ToMap() map[K]V {
	result := make(map[K]V, len(m.items))
	for key, item := range m.items {
		result[key] = item.value
	}
	return result
}

func (m *_map[K, V]) addBeforeInternal(key K, value V, beforeItem *item[K, V]) {
	m.Remove(key)
	newItem := &item[K, V]{
		key:   key,
		value: value,
	}
	m.items[key] = newItem
	newItem.next = beforeItem
	if beforeItem.prev == nil {
		beforeItem.prev = newItem
		m.firstItem = newItem
		return
	}
	beforeItem.prev.next = newItem
	newItem.prev = beforeItem.prev
	beforeItem.prev = newItem
}

func (m *_map[K, V]) addAfterInternal(key K, value V, afterItem *item[K, V]) {
	m.Remove(key)
	newItem := &item[K, V]{
		key:   key,
		value: value,
	}
	m.items[key] = newItem
	newItem.prev = afterItem
	if afterItem.next == nil {
		afterItem.next = newItem
		m.lastItem = newItem
		return
	}
	afterItem.next.prev = newItem
	newItem.next = afterItem.next
	afterItem.next = newItem
}
