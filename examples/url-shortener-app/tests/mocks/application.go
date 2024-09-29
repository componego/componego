package mocks

import (
	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/internal/application"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/test-server"
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
	components, err := a.Application.ApplicationComponents()
	// Add the component that provides access to test server dependency.
	components = append(components, test_server.NewComponent())
	return components, err
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *ApplicationMock) ApplicationConfigInit(_ componego.ApplicationMode, _ any) (map[string]any, error) {
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
