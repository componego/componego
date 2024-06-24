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
	"testing"

	"github.com/componego/componego/impl/processors"
	"github.com/componego/componego/internal/testing/require"
	"github.com/componego/componego/libs/type-cast"
)

func TestToBool(t *testing.T) {
	processor := processors.ToBool()
	for _, testCase := range getTypeCastCases() {
		expectedValue, expectedErr := type_cast.ToBool(testCase)
		actualValue, actualErr := processor.ProcessData(testCase)
		if expectedErr == nil {
			require.NoError(t, actualErr)
		} else {
			require.EqualError(t, actualErr, expectedErr.Error())
		}
		require.NotPanics(t, func() {
			require.Equal(t, expectedValue, actualValue.(bool))
		})
	}
}

func TestIsBool(t *testing.T) {
	processor := processors.IsBool()
	for _, testCase := range getTypeCastCases() {
		_, isBool := testCase.(bool)
		expectedValue := testCase
		actualValue, actualErr := processor.ProcessData(testCase)
		if isBool {
			require.NoError(t, actualErr)
			require.Equal(t, expectedValue, actualValue)
		} else {
			require.Error(t, actualErr)
			require.Nil(t, actualValue)
		}
	}
}

func TestToInt64(t *testing.T) {
	processor := processors.ToInt64()
	for _, testCase := range getTypeCastCases() {
		expectedValue, expectedErr := type_cast.ToInt64(testCase)
		actualValue, actualErr := processor.ProcessData(testCase)
		if expectedErr == nil {
			require.NoError(t, actualErr)
		} else {
			require.EqualError(t, actualErr, expectedErr.Error())
		}
		require.NotPanics(t, func() {
			require.Equal(t, expectedValue, actualValue.(int64))
		})
	}
}

func TestToFloat64(t *testing.T) {
	processor := processors.ToFloat64()
	for _, testCase := range getTypeCastCases() {
		expectedValue, expectedErr := type_cast.ToFloat64(testCase)
		actualValue, actualErr := processor.ProcessData(testCase)
		if expectedErr == nil {
			require.NoError(t, actualErr)
		} else {
			require.EqualError(t, actualErr, expectedErr.Error())
		}
		require.NotPanics(t, func() {
			require.Equal(t, expectedValue, actualValue.(float64))
		})
	}
}

func TestToString(t *testing.T) {
	processor := processors.ToString()
	for _, testCase := range getTypeCastCases() {
		expectedValue, expectedErr := type_cast.ToString(testCase)
		actualValue, actualErr := processor.ProcessData(testCase)
		if expectedErr == nil {
			require.NoError(t, actualErr)
		} else {
			require.EqualError(t, actualErr, expectedErr.Error())
		}
		require.NotPanics(t, func() {
			require.Equal(t, expectedValue, actualValue.(string))
		})
	}
}

func getTypeCastCases() []any {
	// noinspection ALL
	return []any{
		nil,
		true,
		false,
		123,
		0,
		int8(123),
		int8(0),
		int16(123),
		int16(0),
		int32(123),
		int32(0),
		int64(123),
		int64(0),
		uint8(123),
		uint8(0),
		uint16(123),
		uint16(0),
		uint32(123),
		uint32(0),
		uint64(123),
		uint64(0),
		float32(0.0),
		float32(123.123456789),
		float64(0.0),
		float64(123.123456789),
		"true",
		"false",
		"TRUE",
		"FALSE",
		"1",
		"0",
		"-1",
		"",
		"string",
		struct{}{},
		struct{ x int }{x: 123},
		[]int{1, 2, 3},
		[...]string{"1", "2", "3"},
	}
}
