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

package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/componego/componego/internal/utils"
)

type VariableConfig struct {
	Indent     string
	UseNewLine bool
	MaxDepth   int
	MaxItems   int
	Converter  func(value any) any
}

// RenderVariable returns the text representation of variable.
func RenderVariable(instance any, config *VariableConfig) string {
	if config == nil {
		config = &VariableConfig{
			MaxDepth: 3,
		}
	}
	builder := &strings.Builder{}
	r := &renderer{
		writer:  builder,
		ptrSeen: make(map[any]struct{}, 3),
		config:  config,
	}
	if config.UseNewLine {
		r.separator = "\n"
	} else {
		r.separator = " "
	}
	r.process(instance, 0)
	return builder.String()
}

type renderer struct {
	writer    io.Writer
	ptrSeen   map[any]struct{}
	config    *VariableConfig
	separator string
}

func (r *renderer) process(instance any, depth int) {
	if instance == nil {
		utils.Fprint(r.writer, "nil")
		return
	}
	reflectValue := reflect.ValueOf(instance)
	reflectType := reflectValue.Type()
	if depth >= r.config.MaxDepth && r.config.MaxDepth > 0 {
		r.readDepthExceeded(reflectType)
		return
	}
	if instance, ok := instance.(reflect.Type); ok {
		utils.Fprint(r.writer, getTypeName(instance))
		return
	}
	if instance, ok := instance.(reflect.Value); ok {
		utils.Fprint(r.writer, getTypeName(instance.Type()))
		return
	}
	if reflectValue.Kind() == reflect.Pointer {
		pointer := reflectValue.UnsafePointer()
		if _, ok := r.ptrSeen[pointer]; ok {
			utils.Fprintf(r.writer, "%s{[... encountered a cycle]}", getTypeName(reflectType))
			return
		}
		r.ptrSeen[pointer] = struct{}{}
		defer delete(r.ptrSeen, pointer)
	}
	if r.config.Converter != nil {
		instance = r.config.Converter(instance)
	}
	reflectValue = reflect.ValueOf(utils.Indirect(instance))
	r.processByType(reflectType, reflectValue, depth)
}

func (r *renderer) processByType(reflectType reflect.Type, reflectValue reflect.Value, depth int) {
	if depth == 0 && reflectValue.Kind() != reflect.Func {
		utils.Fprint(r.writer, getTypeName(reflectType))
	}
	switch reflectValue.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.String:
		r.simpleType(reflectValue, depth)
		return
	case reflect.Func:
		r.funcType(reflectValue)
		return
	}
	if depth+1 >= r.config.MaxDepth && r.config.MaxDepth > 0 {
		r.readDepthExceeded(reflectType)
		return
	}
	switch reflectValue.Kind() {
	case reflect.Array, reflect.Slice:
		r.sliceType(reflectValue, depth)
	case reflect.Map:
		r.mapType(reflectValue, depth)
	case reflect.Struct:
		r.structType(reflectValue, depth)
	default:
		if depth > 0 {
			utils.Fprint(r.writer, getTypeName(reflectType))
		}
		utils.Fprint(r.writer, "{[... is not supported type]}")
	}
}

func (r *renderer) readDepthExceeded(reflectType reflect.Type) {
	utils.Fprintf(r.writer, "%s{[... read depth exceeded]}", getTypeName(reflectType))
}

func (r *renderer) tabs(count int) {
	if r.config.Indent == "" {
		return
	}
	for i := 0; i < count; i++ {
		utils.Fprint(r.writer, r.config.Indent)
	}
}

func (r *renderer) simpleType(reflectValue reflect.Value, depth int) {
	if depth == 0 {
		utils.Fprint(r.writer, "{ ")
	}
	utils.Fprint(r.writer, toJSON(reflectValue.Interface(), r.config.Indent))
	if depth == 0 {
		utils.Fprint(r.writer, " }")
	}
}

func (r *renderer) funcType(reflectValue reflect.Value) {
	declaration := getTypeName(reflectValue.Type())
	funcForPC := runtime.FuncForPC(reflectValue.Pointer())
	if funcForPC != nil {
		declaration = strings.Replace(declaration, "func", "func "+funcForPC.Name(), 1)
	}
	utils.Fprint(r.writer, declaration)
}

func (r *renderer) sliceType(reflectValue reflect.Value, depth int) {
	utils.Fprint(r.writer, "{")
	sliceLen := reflectValue.Len()
	if sliceLen == 0 {
		utils.Fprint(r.writer, "}")
		return
	}
	depth++
	utils.Fprint(r.writer, r.separator)
	for i := 0; i < sliceLen; i++ {
		r.tabs(depth)
		if i >= r.config.MaxItems && r.config.MaxItems > 0 && i+1 != sliceLen {
			utils.Fprintf(r.writer, "[...and other %d items],%s", sliceLen-r.config.MaxItems, r.separator)
			break
		}
		utils.Fprint(r.writer, i, ": ")
		r.process(reflectValue.Index(i).Interface(), depth)
		utils.Fprint(r.writer, ",", r.separator)
	}
	r.tabs(depth - 1)
	utils.Fprint(r.writer, "}")
}

func (r *renderer) mapType(reflectValue reflect.Value, depth int) {
	utils.Fprint(r.writer, "{")
	reflectKeys := reflectValue.MapKeys()
	if len(reflectKeys) == 0 {
		utils.Fprint(r.writer, "}")
		return
	}
	depth++
	utils.Fprint(r.writer, r.separator)
	keys := make([]*mapKey, len(reflectKeys))
	for i, key := range reflectKeys {
		keys[i] = &mapKey{
			asString: fmt.Sprintf("%v", key),
			asValue:  key,
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].asString < keys[j].asString
	})
	for i, key := range keys {
		r.tabs(depth)
		if i >= r.config.MaxItems && r.config.MaxItems > 0 && i+1 != len(keys) {
			utils.Fprintf(r.writer, "[...and other %d items],%s", len(keys)-r.config.MaxItems, r.separator)
			break
		}
		utils.Fprint(r.writer, toJSON(key.asString, r.config.Indent), ": ")
		r.process(reflectValue.MapIndex(key.asValue).Interface(), depth)
		utils.Fprint(r.writer, ",", r.separator)
	}
	r.tabs(depth - 1)
	utils.Fprint(r.writer, "}")
}

func (r *renderer) structType(reflectValue reflect.Value, depth int) {
	reflectType := reflectValue.Type()
	numField := reflectValue.NumField()
	hasExportedFields := false
	for i := 0; i < numField; i++ {
		if isExportedName(reflectType.Field(i).Name) {
			hasExportedFields = true
			break
		}
	}
	utils.Fprint(r.writer, "{")
	if !hasExportedFields {
		utils.Fprint(r.writer, "[... no public data]}")
		return
	}
	depth++
	utils.Fprint(r.writer, r.separator)
	for i := 0; i < numField; i++ {
		fieldType := reflectType.Field(i)
		name := fieldType.Name
		if !isExportedName(name) {
			continue
		}
		value := reflectValue.Field(i).Interface()
		if tag, ok := fieldType.Tag.Lookup("json"); ok {
			tagOptions := strings.Split(tag, ",")
			if len(tagOptions) > 1 && utils.Contains(tagOptions[1:], "omitempty") && utils.IsEmpty(value) {
				continue
			}
			jsonName := strings.TrimSpace(tagOptions[0])
			if jsonName != "-" {
				name = jsonName
			}
		}
		r.tabs(depth)
		utils.Fprint(r.writer, toJSON(name, ""), ": ")
		r.process(value, depth)
		utils.Fprint(r.writer, ",", r.separator)
	}
	r.tabs(depth - 1)
	utils.Fprint(r.writer, "}")
}

type mapKey struct {
	asString string
	asValue  reflect.Value
}

func toJSON(data any, indent string) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// Disable encoding of special characters.
	// All characters in the result string must be displayed as they are.
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", indent)
	err := encoder.Encode(data)
	if err != nil {
		panic(fmt.Errorf("to beautiful json error: %w", err))
	}
	return string(bytes.TrimRight(buffer.Bytes(), "\n"))
}

func isExportedName(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

func getTypeName(reflectType reflect.Type) string {
	return strings.ReplaceAll(reflectType.String(), "interface {}", "any")
}
