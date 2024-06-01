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

package logger

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// LogData is a function that writes data to the log.
func LogData(writer io.Writer, messages ...any) {
	if writer == nil {
		return
	}
	_, err := fmt.Fprintln(writer, messages...)
	if err != nil {
		panic(err)
	}
}

// ExpectedLogData is a function that returns formatted data.
// You can use it with the LogData function to compare data in tests.
func ExpectedLogData(lines ...any) string {
	var result strings.Builder
	for _, line := range lines {
		if line == nil {
			result.WriteString(fmt.Sprintln(line))
			continue
		}
		reflectValue := reflect.ValueOf(line)
		if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
			result.WriteString(fmt.Sprintln(line))
			continue
		}
		valueLen := reflectValue.Len()
		messages := make([]any, valueLen)
		for i := 0; i < valueLen; i++ {
			messages[i] = reflectValue.Index(i).Interface()
		}
		result.WriteString(fmt.Sprintln(messages...))
	}
	return result.String()
}
