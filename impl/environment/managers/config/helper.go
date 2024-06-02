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

package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/componego/componego"
)

func Get[T any](configKey string, processor componego.Processor, env componego.Environment) (T, error) {
	value, err := env.ConfigProvider().ConfigValue(configKey, processor)
	if err != nil {
		return *new(T), err
	}
	if result, ok := value.(T); ok {
		return result, nil
	}
	return *new(T), fmt.Errorf("could not convert the value for the '%s' to type %T", configKey, *new(T))
}

func GetOrPanic[T any](configKey string, processor componego.Processor, env componego.Environment) T {
	result, err := Get[T](configKey, processor, env)
	if err != nil {
		panic(err)
	}
	return result
}

func ProcessVariables(settings map[string]any) (err error) {
	// ${ENV:VARIABLE_NAME} or ${ENV:VARIABLE_NAME|DEFAULT_VALUE}
	variableRegex := regexp.MustCompile(`\$\{ENV:([a-zA-Z0-9_]+)(\|([a-zA-Z0-9_]+))?}`)
	for key, value := range settings {
		valueAsString, ok := value.(string)
		if !ok {
			// We process nested values since the configuration nesting level can be any.
			valueAsMap, ok := value.(map[string]any)
			if !ok {
				continue
			}
			if err = ProcessVariables(valueAsMap); err != nil {
				return err
			}
			continue
		}
		settings[key] = variableRegex.ReplaceAllStringFunc(valueAsString, func(match string) string {
			matches := variableRegex.FindStringSubmatch(match)
			if envValue := os.Getenv(matches[1]); len(envValue) != 0 {
				return envValue
			} else if len(matches[3]) > 0 {
				return matches[3] // default value
			}
			err = fmt.Errorf("environment variable '%s' not found", matches[1])
			return match
		})
		if err != nil {
			break
		}
	}
	return err
}
