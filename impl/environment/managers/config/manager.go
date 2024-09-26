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

package config

import (
	"strings"

	"github.com/componego/componego"
	"github.com/componego/componego/libs/xerrors"
)

const (
	// delimiter is a constant which specifies which delimiter is used in configuration keys.
	delimiter = "."
)

var (
	ErrConfigManager = xerrors.New("error inside config manager", "E0310")
	ErrConfigInit    = ErrConfigManager.WithMessage("config init error", "E0311")
	ErrConfigGet     = ErrConfigManager.WithMessage("config get error", "E0312")
	ErrValueNotFound = ErrConfigGet.WithMessage("config value not found", "E0313")
)

type manager struct {
	env          componego.Environment
	parsedConfig map[string]any
}

func NewManager() (componego.ConfigProvider, func(componego.Environment, map[string]any) error) {
	m := &manager{}
	return m, m.initialize
}

func (m *manager) ConfigValue(configKey string, processor componego.Processor) (any, error) {
	value, ok := m.extractValue(configKey)
	if processor == nil {
		if ok {
			return value, nil
		}
		return nil, ErrValueNotFound.WithOptions("E0314",
			xerrors.NewOption("componego:config:key", configKey),
		)
	}
	// Injecting a dependency into an object before calling the object's methods.
	err := m.env.DependencyInvoker().PopulateFields(processor)
	if err == nil {
		value, err = processor.ProcessData(value)
		if err == nil {
			return value, nil
		}
	}
	return nil, ErrConfigGet.WithError(err, "E0315",
		xerrors.NewOption("componego:config:key", configKey),
	)
}

func (m *manager) extractValue(configKey string) (any, bool) {
	if value, ok := m.parsedConfig[configKey]; ok {
		return value, true
	}
	keys := strings.Split(configKey, delimiter)
	var configValue any = m.parsedConfig
	for _, key := range keys {
		switch parsedConfig := configValue.(type) {
		case map[string]any:
			if value, ok := parsedConfig[key]; ok {
				configValue = value
			} else {
				return nil, false
			}
		}
	}
	return configValue, true
}

func (m *manager) initialize(env componego.Environment, parsedConfig map[string]any) error {
	m.env = env
	m.parsedConfig = parsedConfig
	return nil
}

func ParseConfig(env componego.Environment, options any) (map[string]any, error) {
	if app, ok := env.Application().(componego.ApplicationConfigInit); ok {
		parsedConfig, err := app.ApplicationConfigInit(env.ApplicationMode(), options)
		if err != nil {
			return nil, ErrConfigInit.WithError(err, "E0316")
		}
		return parsedConfig, nil
	}
	return nil, nil
}
