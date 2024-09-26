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

package processors

import (
	"fmt"

	"github.com/componego/componego"
)

type processor struct {
	handler func(value any) (any, error)
}

func New(handler func(value any) (any, error)) componego.Processor {
	return &processor{
		handler: handler,
	}
}

func (p *processor) ProcessData(value any) (any, error) {
	return p.handler(value)
}

type multiProcessor struct {
	di         componego.DependencyInvoker `componego:"inject"`
	processors []componego.Processor
}

func Multi(processors ...componego.Processor) componego.Processor {
	return &multiProcessor{
		processors: processors,
	}
}

func (m *multiProcessor) ProcessData(value any) (any, error) {
	for _, item := range m.processors {
		if err := m.di.PopulateFields(item); err != nil {
			return nil, err
		} else if value, err = item.ProcessData(value); err != nil {
			return nil, err
		}
	}
	return value, nil
}

func DefaultValue(value any) componego.Processor {
	return New(func(prevValue any) (any, error) {
		if prevValue == nil {
			return value, nil
		}
		return prevValue, nil
	})
}

func IsRequired() componego.Processor {
	return New(func(value any) (any, error) {
		if value == nil {
			return nil, fmt.Errorf("the value is required")
		}
		return value, nil
	})
}
