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

package runner

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/impl/runner/unhandled-errors"
	"github.com/componego/componego/internal/developer"
	"github.com/componego/componego/internal/system"
	"github.com/componego/componego/internal/utils"
)

// RunWithContext runs the application with context and returns the exit code.
func RunWithContext(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) int {
	d := driver.New(nil)
	exitCode, err := d.RunApplication(ctx, app, appMode)
	if err != nil {
		// Here we display all errors that were not processed.
		utils.Fprint(system.Stderr, unhandled_errors.ToString(err, appMode, unhandled_errors.GetHandlers()))
	}
	return exitCode
}

// Run runs the application and returns the exit code.
func Run(app componego.Application, appMode componego.ApplicationMode) int {
	return RunWithContext(context.Background(), app, appMode)
}

// RunAndExit runs the application and exits the program after stopping the application.
func RunAndExit(app componego.Application, appMode componego.ApplicationMode) {
	exitCode := Run(app, appMode)
	exit(exitCode, appMode)
}

// RunGracefullyAndExit runs the application and stops it gracefully.
func RunGracefullyAndExit(app componego.Application, appMode componego.ApplicationMode) {
	cancelableCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx() // This will cancel the context if panic occurs.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-interruptChan:
		case <-cancelableCtx.Done():
			signal.Stop(interruptChan)
		}
		cancelCtx()
	}()
	exitCode := RunWithContext(cancelableCtx, app, appMode)
	cancelCtx()       // This will cancel the context unless panic occurs.
	runtime.Gosched() // We switch the runtime so that waiting goroutines can complete their work.
	exit(exitCode, appMode)
}

func exit(exitCode int, appMode componego.ApplicationMode) {
	if appMode == componego.DeveloperMode && system.NumGoroutineBeforeExit() > 1 {
		// In any case, all goroutines will be terminated after exiting the application, but we will show this message.
		developer.Warning(system.Stdout, "The application was stopped, but goroutines were still running.")
		developer.Warning(system.Stdout, "So it may not be stopped correctly. In some cases, this notification may be false.")
		developer.Warning(system.Stdout, "Read more here https://componego.github.io/warnings/goroutine-leak")
	}
	// Make sure you call this function in the root goroutine to ensure the program exits correctly.
	system.Exit(exitCode)
}
