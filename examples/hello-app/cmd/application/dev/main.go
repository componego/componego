package main

import (
	"github.com/componego/componego"
	"github.com/componego/componego/impl/runner"

	"github.com/componego/componego/examples/hello-app/internal/application"
)

func main() {
	runner.RunAndExit(application.New(), componego.DeveloperMode)
}
