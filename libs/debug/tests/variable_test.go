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
	"reflect"
	"strings"
	"testing"

	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/debug"
)

type testStruct1 struct {
	A float64 `json:"a"`
	B []any
	c bool
	D *testStruct2
}

type testStruct2 struct {
	A bool
	B []string
}

func TestRenderVariable(t *testing.T) {
	require.Equal(t, "nil", debug.RenderVariable(nil, nil))
	require.Equal(t, "int{ 1 }", debug.RenderVariable(1, nil))
	require.Equal(t, `string{ "1" }`, debug.RenderVariable("1", nil))
	require.Equal(t, "int", debug.RenderVariable(reflect.TypeOf(1), nil))
	require.Equal(t, "int", debug.RenderVariable(reflect.ValueOf(1), nil))
	require.Regexp(t, `func \S+\(float64, \.\.\.string\) error`, debug.RenderVariable(func(_ float64, _ ...string) error { return nil }, nil))
	require.Equal(t, "struct {}{[... no public data]}", debug.RenderVariable(struct{}{}, nil))
	structObj := &testStruct1{
		A: 0.1,
		B: []any{"1", 2, false},
		c: true,
		D: &testStruct2{
			A: false,
			B: []string{"text"},
		},
	}
	structOutput := `*tests.testStruct1{ "a": 0.1, "B": { 0: "1", 1: 2, 2: false, }, "D": { "A": false, "B": []string{[... read depth exceeded]}, }, }`
	require.Equal(t, structOutput, debug.RenderVariable(structObj, nil))
	structOutput = `*tests.testStruct1{
  "a": 0.1,
  "B": {
    0: "1",
    1: 2,
    2: false,
  },
  "D": {
    "A": false,
    "B": []string{[... read depth exceeded]},
  },
}`
	require.Equal(t, structOutput, debug.RenderVariable(structObj, &debug.VariableConfig{
		Indent:     "  ",
		UseNewLine: true,
		MaxDepth:   3,
		MaxItems:   10,
	}))
	structOutput = strings.ReplaceAll(strings.ReplaceAll(structOutput, "  ", "@@"), "\n", " ")
	require.Equal(t, structOutput, debug.RenderVariable(structObj, &debug.VariableConfig{
		Indent:     "@@",
		UseNewLine: false,
		MaxDepth:   3,
		MaxItems:   10,
	}))
	structOutput = `*tests.testStruct1{
--"a": 0.1,
--"B": []any{[... read depth exceeded]},
--"D": *tests.testStruct2{[... read depth exceeded]},
}`
	require.Equal(t, structOutput, debug.RenderVariable(structObj, &debug.VariableConfig{
		Indent:     "--",
		UseNewLine: true,
		MaxDepth:   2,
		MaxItems:   2,
	}))
	structOutput = `*tests.testStruct1{
  "a": 0.1,
  "B": {
    0: "1",
    [...and other 2 items],
  },
  "D": {
    "A": false,
    "B": {
      0: "text",
    },
  },
}`
	require.Equal(t, structOutput, debug.RenderVariable(structObj, &debug.VariableConfig{
		Indent:     "  ",
		UseNewLine: true,
		MaxDepth:   4,
		MaxItems:   1,
	}))
}
