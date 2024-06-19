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
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// IsErrorType returns true if the type is an error.
func IsErrorType(reflectType reflect.Type) bool {
	return reflectType.Implements(errorType)
}

func Indirect(instance any) any {
	if instance == nil {
		return nil
	}
	if reflectType := reflect.TypeOf(instance); reflectType.Kind() != reflect.Pointer {
		return instance
	}
	reflectValue := reflect.ValueOf(instance)
	for reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue.Interface()
}

func IsEmpty(instance any) bool {
	if instance == nil {
		return true
	}
	reflectValue := reflect.ValueOf(instance)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice:
		return reflectValue.Len() == 0
	case reflect.Pointer:
		if reflectValue.IsNil() {
			return true
		}
		return IsEmpty(reflectValue.Elem().Interface())
	}
	return reflect.DeepEqual(instance, reflect.Zero(reflectValue.Type()).Interface())
}
