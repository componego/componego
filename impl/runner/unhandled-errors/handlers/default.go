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
	errs := xerrors.UnwrapAll(err)
	xErrs := make([]xerrors.XError, 0, len(errs))
	duplicateCodes := make(map[string]int, len(errs))
	for _, err := range errs {
		// noinspection ALL
		xErr, ok := err.(xerrors.XError) //nolint:errorlint
		if !ok || len(xErr.ErrorOptions()) == 0 {
			continue
		}
		xErrs = append(xErrs, xErr)
		duplicateCodes[xErr.ErrorCode()]++
	}
	if len(xErrs) == 0 {
		return
	}
	prepareDuplicates(duplicateCodes)
	utils.Fprintln(theme.GetWriter("blockName"), "Options:")
	for i, xErr := range xErrs {
		utils.Fprint(theme.GetWriter("default"), i+1, ". ")
		errCode := xErr.ErrorCode()
		if errCode == "" {
			utils.Fprint(theme.GetWriter("default"), "<empty>")
		} else {
			utils.Fprint(theme.GetWriter("errorText"), errCode)
		}
		if _, ok := duplicateCodes[errCode]; ok {
			utils.Fprint(theme.GetWriter("default"), " (", duplicateCodes[errCode], ")")
			duplicateCodes[errCode]++
		}
		utils.Fprintln(theme.GetWriter("default"), ":")
		renderOptions(theme, xErr.ErrorOptions())
	}
}

func renderOptions(theme color.Theme, options []xerrors.Option) {
	duplicateKeys := make(map[string]int, len(options))
	for _, option := range options {
		duplicateKeys[option.Key()]++
	}
	prepareDuplicates(duplicateKeys)
	for _, option := range options {
		optionKey := option.Key()
		utils.Fprint(theme.GetWriter("default"), " * ")
		utils.Fprint(theme.GetWriter("optionKey"), optionKey)
		if _, ok := duplicateKeys[optionKey]; ok {
			utils.Fprint(theme.GetWriter("default"), " (", duplicateKeys[optionKey], ")")
			duplicateKeys[optionKey]++
		}
		utils.Fprint(theme.GetWriter("default"), " => ")
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

func prepareDuplicates(values map[string]int) {
	for key, value := range values {
		if value <= 1 {
			delete(values, key)
		} else {
			values[key] = 0
		}
	}
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
