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

package developer

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func NumGoroutineBeforeExit() int {
	numGoroutine := runtime.NumGoroutine()
	if numGoroutine == 1 {
		return numGoroutine
	}
	// Signals start a goroutine that never stops.
	// We make sure that this goroutine does not influence the result.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGTERM)
	signal.Stop(interruptChan)
	return runtime.NumGoroutine() - 1
}
