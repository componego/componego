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

package processors

import (
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/libs/type-cast"
)

// ToBool converts the value to boolean.
func ToBool() componego.Processor {
	return New(func(value any) (any, error) {
		return type_cast.ToBool(value)
	})
}

// IsBool checks whether a value is a boolean value.
func IsBool() componego.Processor {
	return New(func(value any) (any, error) {
		if _, ok := value.(bool); ok {
			return value, nil
		}
		return nil, fmt.Errorf("the value is not a boolean")
	})
}

// ToInt64 converts the value to int64.
func ToInt64() componego.Processor {
	return New(func(value any) (any, error) {
		return type_cast.ToInt64(value)
	})
}

// ToFloat64 converts the value to float64.
func ToFloat64() componego.Processor {
	return New(func(value any) (any, error) {
		return type_cast.ToFloat64(value)
	})
}

// ToString converts the value to string.
func ToString() componego.Processor {
	return New(func(value any) (any, error) {
		return type_cast.ToString(value)
	})
}
