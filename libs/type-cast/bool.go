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

package type_cast

import (
	"fmt"
	"strconv"

	"github.com/componego/componego/internal/utils"
)

func ToBool(value any) (bool, error) {
	value = utils.Indirect(value)
	switch castedValue := value.(type) {
	case nil:
		return false, nil
	case bool:
		return castedValue, nil
	case int:
		return castedValue != 0, nil
	case int8:
		return castedValue != 0, nil
	case int16:
		return castedValue != 0, nil
	case int32:
		return castedValue != 0, nil
	case int64:
		return castedValue != 0, nil
	case uint:
		return castedValue != 0, nil
	case uint8:
		return castedValue != 0, nil
	case uint16:
		return castedValue != 0, nil
	case uint32:
		return castedValue != 0, nil
	case uint64:
		return castedValue != 0, nil
	case float32:
		return castedValue != 0, nil
	case float64:
		return castedValue != 0, nil
	case string:
		return strconv.ParseBool(castedValue)
	}
	return false, fmt.Errorf("unable to cast %#v of type %T to bool", value, value)
}
