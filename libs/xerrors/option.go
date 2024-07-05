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

type Option interface {
	Key() string
	Value() any
}

type option struct {
	key   string
	value any
}

func NewOption(key string, value any) Option {
	return &option{
		key:   key,
		value: value,
	}
}

func (o *option) Key() string {
	return o.key
}

func (o *option) Value() any {
	return o.value
}

type callableOption struct {
	key      string
	getValue func() any
}

func NewCallableOption(key string, getValue func() any) Option {
	return &callableOption{
		key:      key,
		getValue: getValue,
	}
}

func (c *callableOption) Key() string {
	return c.key
}

func (c *callableOption) Value() any {
	return c.getValue()
}
