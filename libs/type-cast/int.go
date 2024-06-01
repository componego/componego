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

package type_cast

import (
	"fmt"
	"strconv"

	"github.com/componego/componego/internal/utils"
)

func ToInt64(value any) (int64, error) {
	value = utils.Indirect(value)
	switch castedValue := value.(type) {
	case nil:
		return 0, nil
	case bool:
		if castedValue {
			return 1, nil
		}
		return 0, nil
	case int:
		return int64(castedValue), nil
	case int8:
		return int64(castedValue), nil
	case int16:
		return int64(castedValue), nil
	case int32:
		return int64(castedValue), nil
	case int64:
		return castedValue, nil
	case uint:
		return int64(castedValue), nil
	case uint8:
		return int64(castedValue), nil
	case uint16:
		return int64(castedValue), nil
	case uint32:
		return int64(castedValue), nil
	case uint64:
		return int64(castedValue), nil
	case float32:
		return int64(castedValue), nil
	case float64:
		return int64(castedValue), nil
	case string:
		return strconv.ParseInt(castedValue, 0, 0)
	}
	return 0, fmt.Errorf("unable to cast %#v of type %T to int64", value, value)
}
