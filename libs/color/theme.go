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

package color

import (
	"fmt"
	"io"
)

var (
	themeFactories map[string]themeFactory
)

func init() {
	themeFactories = make(map[string]themeFactory, 5)
}

type themeFactory = func(writer io.Writer) Theme

type Theme interface {
	GetWriter(name string) io.Writer
}

type theme struct {
	defaultWriter io.Writer
	colors        map[string]io.Writer
}

func NewTheme(writer io.Writer, colors map[string][]Color) Theme {
	theme := &theme{
		defaultWriter: writer,
		colors:        make(map[string]io.Writer, len(colors)),
	}
	for name, colors := range colors {
		theme.colors[name] = NewColoredWriter(writer, colors...)
	}
	return theme
}

func (t *theme) GetWriter(name string) io.Writer {
	if writer, ok := t.colors[name]; ok {
		return writer
	}
	return t.defaultWriter
}

func AddTheme(name string, themeFactory themeFactory) {
	themeFactories[name] = themeFactory
}

func GetTheme(name string, writer io.Writer) Theme {
	if factory, ok := themeFactories[name]; ok {
		return factory(writer)
	}
	// We use panic here because we don't want you to handle an error when the theme doesn't exist.
	// The theme should always exist if you use it.
	panic(fmt.Sprintf("theme '%s' not found", name))
}

func HasTheme(name string) bool {
	_, ok := themeFactories[name]
	return ok
}
