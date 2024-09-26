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

func ToString(value any) (string, error) {
	value = utils.Indirect(value)
	switch castedValue := value.(type) {
	case nil:
		return "", nil
	case bool:
		return strconv.FormatBool(castedValue), nil
	case int:
		return strconv.Itoa(castedValue), nil
	case int8:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case int16:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case int32:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case int64:
		return strconv.FormatInt(castedValue, 10), nil
	case uint:
		return strconv.FormatUint(uint64(castedValue), 10), nil
	case uint8:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case uint16:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case uint32:
		return strconv.FormatInt(int64(castedValue), 10), nil
	case uint64:
		return strconv.FormatUint(castedValue, 10), nil
	case float32:
		return strconv.FormatFloat(float64(castedValue), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(castedValue, 'f', -1, 64), nil
	case string:
		return castedValue, nil
	}
	return "", fmt.Errorf("unable to cast %#v of type %T to string", value, value)
}
