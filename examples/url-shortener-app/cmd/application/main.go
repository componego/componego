package main

import (
	"github.com/componego/componego"
	"github.com/componego/componego/impl/runner"

	"github.com/componego/componego/examples/url-shortener-app/internal/application"
)

func main() {
	// This is an entry point for launching the application in production mode.
	runner.RunAndExit(application.New(), componego.ProductionMode)
}
