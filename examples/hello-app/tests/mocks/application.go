package mocks

import (
	"github.com/componego/componego/examples/hello-app/internal/application"
)

type ApplicationMock struct {
	*application.Application
}

func NewApplicationMock() *ApplicationMock {
	return &ApplicationMock{
		Application: application.New(),
	}
}
