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

	"github.com/componego/componego"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/impl/runner/unhandled-errors"
	"github.com/componego/componego/internal/developer"
	"github.com/componego/componego/internal/utils"
)

// Run runs the application.
func Run(app componego.Application, appMode componego.ApplicationMode) int {
	d := driver.New(nil)
	exitCode, err := d.RunApplication(context.Background(), app, appMode)
	if err != nil {
		// Here we display all errors that were not processed.
		utils.Fprint(os.Stderr, unhandled_errors.ToString(err, appMode, unhandled_errors.GetHandlers()))
	}
	return exitCode
}

// RunAndExit runs the application and exits the program after stopping the application.
func RunAndExit(app componego.Application, appMode componego.ApplicationMode) {
	exitCode := Run(app, appMode)
	if appMode == componego.DeveloperMode && developer.NumGoroutineBeforeExit() > 1 {
		// In any case, all goroutines will be terminated after exiting the application, but we will show this message.
		developer.Warning(os.Stdout, "The application was stopped, but goroutines were still running.")
		developer.Warning(os.Stdout, "So it may not be stopped correctly. In some cases, this notification may be false.")
		developer.Warning(os.Stdout, "Read more here https://componego.github.io/warnings/goroutine-leak")
	}
	// Make sure you call this function in the root goroutine to ensure the program exits correctly.
	os.Exit(exitCode)
}
