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

package handlers

import (
	"fmt"
	"io"

	"github.com/componego/componego"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/color"
	"github.com/componego/componego/libs/debug"
	"github.com/componego/componego/libs/xerrors"
)

func DefaultHandler(err error, writer io.Writer, _ componego.ApplicationMode) {
	theme := getUnhandledErrorTheme(writer)
	utils.Fprintln(theme.GetWriter("topMessage"), "> An unhandled error has occurred in the application.")
	utils.Fprintln(theme.GetWriter("blockName"), "Details:")
	utils.Fprintln(theme.GetWriter("errorText"), err)
	errOptions := xerrors.GetAllOptions(err)
	if len(errOptions) == 0 {
		return
	}
	hasDuplicateKeys := false
	keys := make(map[string]int, len(errOptions))
	for _, option := range errOptions {
		key := option.Key()
		keys[key] = keys[key] + 1
		if keys[key] > 1 {
			hasDuplicateKeys = true
		}
	}
	utils.Fprintln(theme.GetWriter("blockName"), "Options:")
	for i, option := range errOptions {
		utils.Fprint(theme.GetWriter("default"), i+1, ". ")
		if hasDuplicateKeys {
			key := option.Key()
			keys[key] = keys[key] - 1
			utils.Fprint(theme.GetWriter("optionKey"), key, " (", keys[key], ") => ")
		} else {
			utils.Fprint(theme.GetWriter("optionKey"), option.Key(), " => ")
		}
		// We render everything in one goroutine to be able to correctly catch an error during rendering error options.
		utils.Fprintln(theme.GetWriter("optionValue"), renderVariable(option.Value()))
	}
}

func renderVariable(value any) string {
	switch value := value.(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	}
	return debug.RenderVariable(value, &debug.VariableConfig{
		Indent:     utils.Indent,
		UseNewLine: true,
		MaxDepth:   3,
		MaxItems:   10,
		Converter: func(value any) any {
			switch obj := value.(type) {
			case componego.Application:
				return struct {
					Name string `json:"applicationName"`
				}{
					Name: obj.ApplicationName(),
				}
			case componego.Component:
				return struct {
					Identifier string `json:"componentIdentifier"`
					Version    string `json:"componentVersion"`
				}{
					Identifier: obj.ComponentIdentifier(),
					Version:    obj.ComponentVersion(),
				}
			}
			return value
		},
	})
}

func getUnhandledErrorTheme(writer io.Writer) color.Theme {
	if color.HasTheme("componego:unhandledError") {
		return color.GetTheme("componego:unhandledError", writer)
	}
	return color.NewTheme(writer, map[string][]color.Color{
		"errorText":   {color.RedBackground, color.WhiteColor},
		"topMessage":  {color.RedBackground, color.WhiteColor},
		"blockName":   {color.UnderlineText},
		"optionKey":   {color.GreenColor},
		"optionValue": {color.CyanColor},
	})
}
