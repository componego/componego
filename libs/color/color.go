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
	"bytes"
	"io"
	"strconv"
)

type Color int

const (
	BlackColor  Color = 30
	BlueColor   Color = 34
	CyanColor   Color = 36
	GrayColor   Color = 37
	GreenColor  Color = 32
	PurpleColor Color = 35
	RedColor    Color = 31
	WhiteColor  Color = 97
	YellowColor Color = 33

	BlackBackground  Color = 40
	BlueBackground   Color = 44
	CyanBackground   Color = 46
	GrayBackground   Color = 47
	GreenBackground  Color = 42
	PurpleBackground Color = 45
	RedBackground    Color = 41
	WhiteBackground  Color = 107
	YellowBackground Color = 43

	BoldText      Color = 1
	UnderlineText Color = 4

	reset Color = 0
)

var (
	_active = false // Colors are disabled by default
)

func SetIsActive(active bool) {
	_active = active
}

func GetIsActive() bool {
	return _active
}

type coloredWriter struct {
	writer io.Writer
	colors []Color
}

func NewColoredWriter(writer io.Writer, colors ...Color) io.Writer {
	return &coloredWriter{
		writer: writer,
		colors: colors,
	}
}

func (s *coloredWriter) Write(data []byte) (int, error) {
	if !_active || len(s.colors) == 0 {
		return s.writer.Write(data)
	}
	data = normalize(data)
	if bytes.IndexByte(data, '\n') < 0 {
		newData := make([]byte, 0, len(data)+(len(s.colors)+1)*5) // 5 is the approximate size of one color in bytes.
		newData = appendColors(newData, data, s.colors)
		return s.writer.Write(newData)
	}
	// In some cases, colors may color a new line if there is a new line.
	// Colors will not be applied to the new line because it will look ugly.
	items := bytes.Split(data, []byte{'\n'})
	newData := make([]byte, 0, len(data)+(len(s.colors)+1)*len(items)*5) // 5 is the approximate size of one color in bytes.
	for i, data := range items {
		if i > 0 {
			newData = append(newData, '\n')
		}
		newData = appendColors(newData, data, s.colors)
	}
	return s.writer.Write(newData)
}

func normalize(data []byte) []byte {
	if bytes.IndexByte(data, '\r') < 0 {
		return data
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte{'\n'})
	if bytes.IndexByte(data, '\r') < 0 {
		return data
	}
	return bytes.ReplaceAll(data, []byte{'\r'}, []byte{'\n'})
}

func appendColors(newData, data []byte, colors []Color) []byte {
	if len(data) == 0 {
		return newData
	}
	for _, color := range colors {
		newData = colorToBytes(newData, color)
	}
	newData = append(newData, data...)
	return colorToBytes(newData, reset)
}

func colorToBytes(data []byte, color Color) []byte {
	data = append(data, "\033["...)
	data = append(data, []byte(strconv.Itoa(int(color)))...)
	data = append(data, 'm')
	return data
}
