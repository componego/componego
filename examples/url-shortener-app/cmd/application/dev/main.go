package main

import (
	"github.com/componego/componego"
	"github.com/componego/componego/impl/runner"
	"github.com/componego/componego/libs/color"

	"github.com/componego/componego/examples/url-shortener-app/internal/application"
)

func main() {
	// This enables color text output.
	color.SetIsActive(true)
	// This is an entry point for launching the application in developer mode.
	runner.RunGracefullyAndExit(application.New(), componego.DeveloperMode)
}
