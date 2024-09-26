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

package componego

type Dependency any

// DependencyInvoker is an interface that describes the functions that invoke dependencies for the application.
type DependencyInvoker interface {
	// Invoke uses dependencies.
	// You can only pass a function as an argument.
	// If the passed function returns an error, then it will be returned.
	Invoke(function any) (any, error)
	// Populate is similar to Invoke, but it can only take objects that are a reference.
	// The passed object will be filled with dependencies.
	Populate(target any) error
	// PopulateFields fills the struct fields passed as an argument with dependencies.
	// Fields must have a special tag: `componego:"inject"`.
	// Fields without this tag are ignored.
	PopulateFields(target any) error
}
