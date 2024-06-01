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
	"fmt"
	"reflect"

	"github.com/componego/componego"
)

func GenerateTestFactories(countFactories int, countReturnTypes int) []componego.Dependency {
	result := make([]componego.Dependency, countFactories)
	for i := 0; i < countFactories; i++ {
		returnTypes := make([]reflect.Type, countReturnTypes)
		returnValues := make([]reflect.Value, countReturnTypes)
		for j := 0; j < countReturnTypes; j++ {
			returnType := reflect.StructOf([]reflect.StructField{
				{
					Name: fmt.Sprintf("Field_%d_%d", i, j),
					Type: reflect.TypeOf(float64(0)),
				},
			})
			returnTypes[j] = reflect.PointerTo(returnType)
			returnValues[j] = reflect.Zero(returnTypes[j])
		}
		funcType := reflect.FuncOf([]reflect.Type{}, returnTypes, false)
		result[i] = reflect.MakeFunc(funcType, func(_ []reflect.Value) []reflect.Value {
			return returnValues
		}).Interface()
	}
	return result
}
