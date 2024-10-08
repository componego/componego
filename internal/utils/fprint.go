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

import (
	"fmt"
	"io"
)

const Indent = "  "

// noinspection SpellCheckingInspection
func Fprint(writer io.Writer, args ...any) {
	_, err := fmt.Fprint(writer, args...)
	if err != nil {
		panic(err)
	}
}

// noinspection SpellCheckingInspection
func Fprintln(writer io.Writer, args ...any) {
	// fmt.Fprintln adds unnecessary spaces between arguments.
	// This method brings the output to a single format without spaces.
	if len(args) > 0 {
		Fprint(writer, args...)
	}
	Fprint(writer, "\n")
}

func Fprintf(writer io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(writer, format, args...)
	if err != nil {
		panic(err)
	}
}
