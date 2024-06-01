package mocks

import (
	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/internal/application"
)

type ApplicationMock struct {
	*application.Application
}

func NewApplicationMock() *ApplicationMock {
	return &ApplicationMock{
		Application: application.New(),
	}
}

// ApplicationComponents belongs to interface componego.ApplicationComponents.
func (a *ApplicationMock) ApplicationComponents() ([]componego.Component, error) {
	// In this function we call the original application method and can modify it in our mock.
	// For example, you can add or remove some returned components or completely replace the list with your components.
	// You can also create a similar method for other application methods to replace other parts of the application in the mock.
	components, err := a.Application.ApplicationComponents()
	// ...
	return components, err
}

func (a *ApplicationMock) ApplicationDependencies() ([]componego.Dependency, error) {
	dependencies, err := a.Application.ApplicationDependencies()
	// Of course, you must rewrite entities according to the rewrite rules.
	// These rules are written in the documentation.
	dependencies = append(dependencies, NewRedirectRepositoryMock)
	return dependencies, err
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *ApplicationMock) ApplicationConfigInit(_ componego.ApplicationMode) (map[string]any, error) {
	return map[string]any{
		"databases.main-storage.driver": "db-driver-mock",
		"databases.main-storage.source": "...",
	}, nil
}

var (
	_ componego.Application             = (*ApplicationMock)(nil)
	_ componego.ApplicationComponents   = (*ApplicationMock)(nil)
	_ componego.ApplicationDependencies = (*ApplicationMock)(nil)
	_ componego.ApplicationConfigInit   = (*ApplicationMock)(nil)
	_ componego.ApplicationErrorHandler = (*ApplicationMock)(nil)
)
