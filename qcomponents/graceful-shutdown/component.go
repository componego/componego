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

package graceful_shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/componego/componego"
)

type Component struct {
}

func NewComponent() *Component {
	return &Component{}
}

// ComponentIdentifier belongs to interface componego.Component.
func (c *Component) ComponentIdentifier() string {
	return "componego:graceful-shutdown"
}

// ComponentVersion belongs to interface componego.Component.
func (c *Component) ComponentVersion() string {
	return "0.0.1"
}

// ComponentInit belongs to interface componego.ComponentInit.
func (c *Component) ComponentInit(env componego.Environment) error {
	cancelableCtx, cancelCtx := context.WithCancel(env.GetContext())
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
	return env.SetContext(cancelableCtx)
}

var (
	_ componego.Component     = (*Component)(nil)
	_ componego.ComponentInit = (*Component)(nil)
)
