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

package tests

import (
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency/container"
)

func TestDependencyContainer(t *testing.T) {
	DependencyContainerTester[*testing.T](t, func() (container.Container, func([]componego.Dependency) (func() error, error)) {
		return container.New(5)
	})
}

func BenchmarkDependencyContainerInitialize(b *testing.B) {
	factories := GenerateTestFactories(1000, 5)
	b.Run("dependency container initialize", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, initializer := container.New(len(factories))
			if _, err := initializer(factories); err != nil {
				b.Fatal(err)
			}
		}
	})
}
