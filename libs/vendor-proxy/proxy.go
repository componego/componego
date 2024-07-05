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

package vendor_proxy

import (
	"reflect"
	"sync"

	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/xerrors"
)

var (
	ErrVendorProxy           = xerrors.New("vendor proxy error", "E0810")
	ErrFunctionExists        = ErrVendorProxy.WithMessage("such a function already exists", "E0811")
	ErrFunctionNotExists     = ErrVendorProxy.WithMessage("function does not exist", "E0812")
	ErrInvalidArgument       = ErrVendorProxy.WithMessage("function argument is invalid", "E0813")
	ErrArgumentsNotValidated = ErrVendorProxy.WithMessage("arguments are not validated", "E0814")
	ErrWrongCountArguments   = ErrVendorProxy.WithMessage("wrong count of arguments", "E0815")

	instances map[string]Proxy
	mutex     sync.Mutex
	nilTypes  = []reflect.Kind{
		reflect.Array,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.Struct,
	}
)

type contextFunction = func(_ Context, _ ...any) (any, error)

func init() {
	instances = make(map[string]Proxy, 5)
	mutex = sync.Mutex{}
}

type Proxy interface {
	AddFunction(name string, function any) error
	CallFunction(name string, args ...any) (any, error)
}

type Context interface {
	Validated()
	ValidationFailed(message string)
}

type proxy struct {
	mutex            sync.RWMutex
	reflectFunctions map[string]*reflectFunction
	contextFunctions map[string]contextFunction
	name             string
}

func Get(name string) Proxy {
	mutex.Lock()
	instance, ok := instances[name]
	if !ok {
		instance = &proxy{
			mutex:            sync.RWMutex{},
			reflectFunctions: map[string]*reflectFunction{},
			contextFunctions: map[string]contextFunction{},
			name:             name,
		}
		instances[name] = instance
	}
	mutex.Unlock()
	return instance
}

func Has(name string) bool {
	mutex.Lock()
	_, ok := instances[name]
	mutex.Unlock()
	return ok
}

func (p *proxy) AddFunction(name string, function any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if function == nil {
		return p.getError(ErrInvalidArgument, name)
	}
	if function, ok := function.(contextFunction); ok {
		if _, ok := p.contextFunctions[name]; ok {
			return p.getError(ErrFunctionExists, name)
		}
		p.contextFunctions[name] = function
		return nil
	}
	reflectType := reflect.TypeOf(function)
	if reflectType.Kind() != reflect.Func {
		return p.getError(ErrInvalidArgument, name)
	}
	if _, ok := p.reflectFunctions[name]; ok {
		return p.getError(ErrFunctionExists, name)
	}
	p.reflectFunctions[name] = &reflectFunction{
		reflectValue: reflect.ValueOf(function),
		reflectType:  reflectType,
	}
	return nil
}

func (p *proxy) CallFunction(name string, args ...any) (result any, err error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if function := p.contextFunctions[name]; function != nil {
		result, err = p.contextCall(name, function, args...)
	} else if function := p.reflectFunctions[name]; function != nil {
		result, err = p.reflectCall(name, function, args...)
	} else {
		err = p.getError(ErrFunctionNotExists, name)
	}
	return
}

func (p *proxy) getError(err xerrors.XError, functionName string) error {
	return err.WithOptions("E0816",
		xerrors.NewOption("vendorProxy:name", p.name),
		xerrors.NewOption("vendorProxy:functionName", functionName),
	)
}

func (p *proxy) contextCall(functionName string, function contextFunction, args ...any) (result any, err error) {
	ctx := &context{}
	defer func() {
		r := recover()
		if r == nil {
			if !ctx.validated {
				// make sure you call ctx.Validated() in your code after casting the function arguments.
				err = p.getError(ErrArgumentsNotValidated, functionName)
			}
		} else if ctx.validationFailed || !ctx.validated {
			err = p.getError(
				ErrInvalidArgument.WithOptions("E0817",
					xerrors.NewOption("vendorProxy:panic", r),
				),
				functionName,
			)
		} else {
			panic(r)
		}
	}()
	result, err = function(ctx, args...)
	// function return variables using context must match the behavior of function return values using reflect.
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *proxy) reflectCall(functionName string, function *reflectFunction, args ...any) (any, error) {
	reflectType := function.reflectType
	argLen := len(args)
	numIn := reflectType.NumIn()
	var input []reflect.Value
	if reflectType.IsVariadic() {
		if argLen == numIn-1 {
			input = make([]reflect.Value, argLen)
		} else if numIn != argLen {
			return nil, p.getError(ErrWrongCountArguments, functionName)
		} else {
			argLen--
			if args[argLen] == nil {
				return nil, p.getError(ErrInvalidArgument, functionName)
			}
			lastItem := reflect.ValueOf(args[argLen])
			if lastItem.Kind() != reflect.Slice && lastItem.Kind() != reflect.Array {
				return nil, p.getError(ErrInvalidArgument, functionName)
			}
			lastItemLen := lastItem.Len()
			input = make([]reflect.Value, argLen+lastItemLen)
			itemType := reflectType.In(argLen).Elem()
			for i := 0; i < lastItemLen; i++ {
				item := lastItem.Index(i)
				if !item.Type().AssignableTo(itemType) {
					return nil, p.getError(ErrInvalidArgument, functionName)
				}
				input[i+argLen] = item
			}
		}
	} else if numIn != argLen {
		return nil, p.getError(ErrWrongCountArguments, functionName)
	} else {
		input = make([]reflect.Value, argLen)
	}
	for i := 0; i < argLen; i++ {
		item, err := p.getArgItem(functionName, args[i], reflectType.In(i))
		if err != nil {
			return nil, err
		}
		input[i] = item
	}
	output := function.reflectValue.Call(input)
	numOut := len(output)
	if numOut > 0 && utils.IsErrorType(output[numOut-1].Type()) {
		if err := output[numOut-1].Interface(); err != nil {
			return nil, err.(error)
		}
		numOut--
	}
	if numOut == 0 {
		return nil, nil
	}
	if numOut == 1 {
		return output[0].Interface(), nil
	}
	result := make([]any, numOut)
	for i := 0; i < numOut; i++ {
		result[i] = output[i].Interface()
	}
	return result, nil
}

func (p *proxy) getArgItem(functionName string, item any, requireType reflect.Type) (reflect.Value, error) {
	if item == nil {
		for _, nilType := range nilTypes {
			if requireType.Kind() == nilType {
				return reflect.New(requireType).Elem(), nil
			}
		}
	} else if reflect.TypeOf(item).AssignableTo(requireType) {
		return reflect.ValueOf(item), nil
	}
	return *new(reflect.Value), p.getError(ErrInvalidArgument, functionName)
}

type reflectFunction struct {
	reflectValue reflect.Value
	reflectType  reflect.Type
}

type context struct {
	validated        bool
	validationFailed bool
}

func (c *context) Validated() {
	c.validated = true
}

func (c *context) ValidationFailed(message string) {
	c.validationFailed = true
	panic(message)
}
