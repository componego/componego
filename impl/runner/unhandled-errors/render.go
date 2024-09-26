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

package unhandled_errors

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/runner/unhandled-errors/handlers"
	"github.com/componego/componego/libs/ordered-map"
)

type handler = func(err error, writer io.Writer, appMode componego.ApplicationMode) bool

func GetHandlers() ordered_map.Map[string, handler] {
	result := ordered_map.New[string, handler](3)
	result.Set("componego:vendor-proxy", handlers.VendorProxyHandler)
	return result
}

func ToString(err error, appMode componego.ApplicationMode, errHandlers ordered_map.Map[string, handler]) (message string) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			message = fmt.Sprintf("panic during rendering the original error: %v\n\nstack: %s\n", r, stack)
		}
	}()
	writer := &strings.Builder{}
	for _, fn := range errHandlers.Values() {
		if fn(err, writer, appMode) {
			return writer.String()
		}
	}
	handlers.DefaultHandler(err, writer, appMode)
	return writer.String()
}
