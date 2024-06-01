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

package container

import (
	"reflect"

	"github.com/componego/componego"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/xerrors"
)

var (
	ErrDependencyContainer  = xerrors.New("error inside dependency container")
	ErrInvalidProvidedType  = ErrDependencyContainer.WithMessage("invalid provided dependency type or factory")
	ErrNilFactory           = ErrInvalidProvidedType.WithMessage("dependency factory is nil")
	ErrVariadicFactory      = ErrInvalidProvidedType.WithMessage("dependency factory has a variable number of arguments")
	ErrNoReturnDependencies = ErrInvalidProvidedType.WithMessage("dependency factory does not return any type")
	ErrSameDependencyType   = ErrInvalidProvidedType.WithMessage("dependency factory returns more than one dependency variable of the same type")
	ErrIncorrectRewrite     = ErrInvalidProvidedType.WithMessage("dependency type substituted incorrectly")
	ErrGettingDependency    = ErrDependencyContainer.WithMessage("error getting dependency for type")
	ErrUndeclaredDependency = ErrGettingDependency.WithMessage("factory accepts an undeclared dependency type")
	ErrCyclicDependencies   = ErrGettingDependency.WithMessage("cycle detected in dependencies")
	ErrNotFoundType         = ErrGettingDependency.WithMessage("dependency of the requested type was not found")
)

type Container interface {
	GetValue(itemType reflect.Type) (reflect.Value, error)
}

type container struct {
	nodes            map[reflect.Type]*node
	initStack        []reflect.Type
	rewritePositions map[int]struct{}
}

func New(approximateSize int) (Container, func([]componego.Dependency) error) {
	c := &container{
		nodes:     make(map[reflect.Type]*node, approximateSize),
		initStack: make([]reflect.Type, 0, 10),
	}
	return c, c.initialize
}

func (c *container) initialize(dependencies []componego.Dependency) (err error) {
	if len(dependencies) == 0 {
		return err
	}
	// List of positions of nodes that have been replaced.
	// As a result, we should not have nodes with these positions.
	c.rewritePositions = map[int]struct{}{}
	for i, item := range dependencies {
		// we increase the index by the size of the previous count of nodes to avoid duplicates in the positions,
		// and add a new node.
		if err = c.addNode(i, item); err != nil {
			return err
		}
	}
	// Here we check that the dependency rewrite was correct.
	if err = c.checkRewrites(); err != nil {
		return err
	}
	c.rewritePositions = nil
	// We initialize all values in one thread without multithreading to avoid race conditions and define cycles correctly.
	for itemType, nodeObj := range c.nodes {
		if err = c.initValue(nodeObj, itemType); err != nil {
			return err
		}
	}
	return err
}

func (c *container) GetValue(itemType reflect.Type) (reflect.Value, error) {
	nodeObj := c.nodes[itemType]
	if nodeObj == nil {
		return *new(reflect.Value), ErrNotFoundType.WithOptions(
			xerrors.NewOption("componego:dependency:container:requestedType", itemType),
		)
	}
	err := c.initValue(nodeObj, itemType)
	return nodeObj.value, err
}

func (c *container) initValue(nodeObj *node, itemType reflect.Type) error {
	// If the factory is missing, then the value has already been initialized.
	if nodeObj.factory == nil {
		return nil
	}
	factoryObj := nodeObj.factory
	if factoryObj.lock {
		return ErrCyclicDependencies.WithOptions(
			xerrors.NewOption("componego:dependency:container:cycle", c.getCyclicDependencies()),
			xerrors.NewOption("componego:dependency:container:factory", factoryObj.value.Type()),
			xerrors.NewOption("componego:dependency:container:requestedType", itemType),
		)
	}
	factoryObj.lock = true
	c.initStack = append(c.initStack, itemType)
	defer func() {
		// We always return to the previous state after the function completes.
		factoryObj.lock = false
		// Pop the current type from the stack since it passed successfully without cycles.
		c.initStack = c.initStack[:len(c.initStack)-1]
	}()
	input := make([]reflect.Value, len(factoryObj.dependencies))
	for i, dependencyType := range factoryObj.dependencies {
		nodeObj = c.nodes[dependencyType]
		// We check that dependency are present.
		if nodeObj == nil {
			return ErrUndeclaredDependency.WithOptions(
				xerrors.NewOption("componego:dependency:container:factory", factoryObj.value.Type()),
				xerrors.NewOption("componego:dependency:container:undeclaredType", dependencyType),
			)
		}
		// We recursively initialize all the values that are needed to initialize the current value.
		if err := c.initValue(nodeObj, dependencyType); err != nil {
			return err
		}
		input[i] = nodeObj.value
	}
	output := factoryObj.value.Call(input)
	outputLen := len(output)
	if factoryObj.hasError {
		// An additional type check is not needed, because we already know that the last value is an error.
		if lastValue := output[outputLen-1].Interface(); lastValue != nil {
			// noinspection ALL
			return ErrGettingDependency.WithError(lastValue.(error),
				xerrors.NewOption("componego:dependency:container:factory", factoryObj.value.Type()),
				xerrors.NewOption("componego:dependency:container:requestedType", itemType),
			)
		}
		outputLen--
	}
	for i := 0; i < outputLen; i++ {
		nodeObj = c.nodes[factoryObj.types[i]]
		nodeObj.value = output[i]
		// Here we mark the current value as initialized.
		nodeObj.factory = nil
	}
	return nil
}

func (c *container) addNode(position int, item componego.Dependency) error {
	if item == nil {
		return ErrNilFactory
	}
	itemType := reflect.TypeOf(item)
	switch itemType.Kind() {
	case reflect.Func: // The element is a dependency factory.
		if itemType.IsVariadic() {
			return ErrVariadicFactory.WithOptions(
				xerrors.NewOption("componego:dependency:container:factory", itemType),
			)
		}
		numIn := itemType.NumIn()
		numOut := itemType.NumOut()
		if numOut == 0 {
			return ErrNoReturnDependencies.WithOptions(
				xerrors.NewOption("componego:dependency:container:factory", itemType),
			)
		}
		factoryObj := &factory{
			value: reflect.ValueOf(item),
		}
		if utils.IsErrorType(itemType.Out(numOut - 1)) { // last value.
			numOut--
			if numOut == 0 {
				return ErrNoReturnDependencies.WithOptions(
					xerrors.NewOption("componego:dependency:container:factory", itemType),
				)
			}
			factoryObj.hasError = true
		}
		factoryObj.dependencies = make([]reflect.Type, numIn)
		factoryObj.types = make([]reflect.Type, numOut)
		// The dependency factory can also accept other types as function arguments.
		for i := 0; i < numIn; i++ {
			factoryObj.dependencies[i] = itemType.In(i)
		}
		// Return types are new dependencies.
		for i := 0; i < numOut; i++ {
			outType := itemType.Out(i)
			if !isAllowedFactoryReturnType(outType) {
				return ErrInvalidProvidedType.WithOptions(
					xerrors.NewOption("componego:dependency:container:factory", itemType),
					xerrors.NewOption("componego:dependency:container:outType", outType),
				)
			}
			// We add the current type to the rewrites if such a type already exists.
			// If the position matches the position of the previous type, then the factory returns 2 or more identical objects.
			if c.addRewriteToCheck(outType) == position {
				return ErrSameDependencyType.WithOptions(
					xerrors.NewOption("componego:dependency:container:factory", itemType),
					xerrors.NewOption("componego:dependency:container:outType", outType),
				)
			}
			factoryObj.types[i] = outType
			// Adds a new type that can be obtained using a factory.
			c.nodes[outType] = &node{
				factory:  factoryObj,
				position: position,
			}
		}
	case reflect.Pointer:
		if itemType.Elem().Kind() != reflect.Struct {
			return ErrInvalidProvidedType.WithOptions(
				xerrors.NewOption("componego:dependency:container:providedType", itemType),
			)
		}
		c.addRewriteToCheck(itemType)
		// There is no need for a factory here because the value is already ready.
		c.nodes[itemType] = &node{
			value:    reflect.ValueOf(item),
			position: position,
		}
	default:
		return ErrInvalidProvidedType.WithOptions(
			xerrors.NewOption("componego:dependency:container:providedType", itemType),
		)
	}
	return nil
}

func (c *container) addRewriteToCheck(itemType reflect.Type) int {
	if prevNode := c.nodes[itemType]; prevNode != nil {
		c.rewritePositions[prevNode.position] = struct{}{}
		return prevNode.position
	}
	// Any non-existent value.
	// We are sure that there are no negative positions.
	return -1
}

func (c *container) checkRewrites() error {
	if len(c.rewritePositions) == 0 {
		return nil
	}
	for _, nodeObj := range c.nodes {
		// Here we check that all dependencies that were replaced are removed.
		if _, ok := c.rewritePositions[nodeObj.position]; !ok {
			continue
		} else if nodeObj.factory == nil {
			return ErrIncorrectRewrite.WithOptions(
				xerrors.NewOption("componego:dependency:container:providedType", nodeObj.value.Type()),
			)
		}
		return ErrIncorrectRewrite.WithOptions(
			xerrors.NewOption("componego:dependency:container:factory", nodeObj.factory.value.Type()),
		)
	}
	return nil
}

func (c *container) getCyclicDependencies() []*CycleItem {
	result := make([]*CycleItem, len(c.initStack))
	for i, itemType := range c.initStack {
		result[i] = &CycleItem{
			ItemType: itemType,
			Factory:  c.nodes[itemType].factory.value.Type(), // Only factories can have cycles.
		}
	}
	return result
}

type CycleItem struct {
	ItemType reflect.Type
	Factory  reflect.Type
}

type factory struct {
	value        reflect.Value
	types        []reflect.Type
	dependencies []reflect.Type
	hasError     bool
	lock         bool
}

type node struct {
	value    reflect.Value
	factory  *factory
	position int
}

func isAllowedFactoryReturnType(reflectType reflect.Type) bool {
	switch reflectType.Kind() {
	case reflect.Interface:
		return true
	case reflect.Func:
		return true
	case reflect.Pointer:
		return reflectType.Elem().Kind() == reflect.Struct
	}
	return false
}
