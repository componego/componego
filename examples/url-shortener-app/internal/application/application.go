package application

import (
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/environment/managers/config"

	"github.com/componego/componego/examples/url-shortener-app/internal/migration"
	"github.com/componego/componego/examples/url-shortener-app/internal/repository"
	"github.com/componego/componego/examples/url-shortener-app/internal/server"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database"

	"github.com/componego/componego/examples/url-shortener-app/third_party/config-reader"
	_ "github.com/componego/componego/examples/url-shortener-app/third_party/db-driver"
)

type Application struct{}

func New() *Application {
	return &Application{}
}

// ApplicationName belongs to interface componego.Application.
func (a *Application) ApplicationName() string {
	return "Url Shortener App v0.0.1"
}

// ApplicationComponents belongs to interface componego.ApplicationComponents.
func (a *Application) ApplicationComponents() ([]componego.Component, error) {
	return []componego.Component{
		// This is the custom component implemented in this example
		// which provides access to the database using a standard database connection interface.
		database.NewComponent(),
	}, nil
}

// ApplicationDependencies belongs to interface componego.ApplicationDependencies.
func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
	return []componego.Dependency{
		// Pay attention to the implementation of the function to understand what dependencies it provides.
		repository.NewRedirectRepository,
	}, nil
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode, _ any) (settings map[string]any, err error) {
	switch appMode {
	case componego.ProductionMode:
		settings, err = config_reader.Read("./config/production.config.json")
	case componego.DeveloperMode:
		settings, err = config_reader.Read("./config/developer.config.json")
	case componego.TestMode:
		settings, err = config_reader.Read("./config/test.config.json")
	default:
		return nil, fmt.Errorf("not supported application mode: %d", appMode)
	}
	if err == nil {
		// If necessary, you can additionally process the settings after reading the configuration.
		// In this case, we add environment variables to the configuration.
		err = config.ProcessVariables(settings)
	}
	return settings, err
}

// ApplicationErrorHandler belongs to interface componego.ApplicationErrorHandler.
func (a *Application) ApplicationErrorHandler(err error, _ componego.ApplicationIO, _ componego.ApplicationMode) error {
	// This method catches all previously unhandled errors.
	// You can process them in some way or return an error.
	// If you return an error, it will be handled at a lower level of the framework.
	return err
}

// ApplicationAction belongs to interface componego.Application.
func (a *Application) ApplicationAction(env componego.Environment, _ any) (int, error) {
	// We always run migrations after the application is initialized.
	if _, err := env.DependencyInvoker().Invoke(migration.Run); err != nil {
		return application.ExitWrapper(err)
	}
	// Start server after the migrations are completed.
	_, err := env.DependencyInvoker().Invoke(server.Run)
	return application.ExitWrapper(err)
}

var (
	_ componego.Application             = (*Application)(nil)
	_ componego.ApplicationComponents   = (*Application)(nil)
	_ componego.ApplicationDependencies = (*Application)(nil)
	_ componego.ApplicationConfigInit   = (*Application)(nil)
	_ componego.ApplicationErrorHandler = (*Application)(nil)
)
