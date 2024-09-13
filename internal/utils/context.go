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

package utils

import (
	"context"
	"reflect"
)

const parentContextField = "Context"

func init() {
	checkContextForUsedGoVersion()
}

func IsParentContext(parent context.Context, child context.Context) (ok bool) {
	if parent == nil || child == nil {
		return false
	}
	for {
		if parent == child {
			return true
		}
		reflectValue := reflect.Indirect(reflect.ValueOf(child))
		if reflectValue.Kind() != reflect.Struct {
			return false
		}
		contextField := reflectValue.FieldByName(parentContextField)
		if !contextField.IsValid() || contextField.IsNil() {
			return false
		}
		child, ok = contextField.Interface().(context.Context)
		if !ok {
			return false
		}
	}
}

func checkContextForUsedGoVersion() {
	errMessage := `
It looks like you are using a GoLang version that is unsupported by your project's framework version.
This incompatibility was not detected during project compilation.
---
You need to update your framework to the latest version.
If this does not solve the problem, please create an issue on GitHub -> https://github.com/componego/componego/issues
We will fix this problem as soon as possible.
`
	parentCtx := context.Background()
	childCtx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	reflectValue := reflect.Indirect(reflect.ValueOf(childCtx))
	if reflectValue.Kind() != reflect.Struct {
		panic(errMessage)
	}
	contextField := reflectValue.FieldByName(parentContextField)
	if !contextField.IsValid() {
		panic(errMessage)
	}
	internalCtx, ok := contextField.Interface().(context.Context)
	if !ok || internalCtx != parentCtx {
		panic(errMessage)
	}
}
