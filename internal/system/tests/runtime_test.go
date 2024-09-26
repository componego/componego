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

package tests

import (
	"sync"
	"testing"

	"github.com/componego/componego/internal/system"
	"github.com/componego/componego/internal/testing/require"
)

func TestNumGoroutineBeforeExit(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		require.True(t, system.NumGoroutineBeforeExit() >= 1)
		numGoroutine := 0
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			numGoroutine = system.NumGoroutineBeforeExit()
		}()
		wg.Wait()
		require.True(t, numGoroutine > 1)
	})
}
