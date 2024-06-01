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

package dependency

import (
	"fmt"

	"github.com/componego/componego"
)

func Get[T any](env componego.Environment) (T, error) {
	value := *new(T)
	err := env.DependencyInvoker().Populate(&value)
	return value, err
}

func GetOrPanic[T any](env componego.Environment) T {
	value, err := Get[T](env)
	if err != nil {
		panic(err)
	}
	return value
}

func Invoke[T any](fn any, env componego.Environment) (T, error) {
	value, err := env.DependencyInvoker().Invoke(fn)
	if err != nil {
		return *new(T), err
	}
	if result, ok := value.(T); ok {
		return result, nil
	}
	return *new(T), fmt.Errorf("сould not convert the returned value to type %T", *new(T))
}

func InvokeOrPanic[T any](fn any, env componego.Environment) T {
	value, err := Invoke[T](fn, env)
	if err != nil {
		panic(err)
	}
	return value
}
