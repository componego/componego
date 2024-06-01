package application

import (
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

// ApplicationName belongs to interface componego.Application.
func (a *Application) ApplicationName() string {
	return "Hello World App v0.0.1"
}

// ApplicationAction belongs to interface componego.Application.
func (a *Application) ApplicationAction(env componego.Environment, _ []string) (int, error) {
	_, err := fmt.Fprintln(env.ApplicationIO().OutputWriter(), "Hello World!")
	return application.ExitWrapper(err)
}

var (
	_ componego.Application = (*Application)(nil)
)
