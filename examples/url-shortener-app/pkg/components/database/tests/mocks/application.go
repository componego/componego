package mocks

import (
	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"

	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database"

	_ "github.com/componego/componego/examples/url-shortener-app/third_party/db-driver"
)

func NewApplicationMock() componego.Application {
	factory := application.NewFactory("Application for Test Database Component")
	factory.SetApplicationComponents(func() ([]componego.Component, error) {
		return []componego.Component{
			database.NewComponent(),
		}, nil
	})
	factory.SetApplicationConfigInit(func(_ componego.ApplicationMode, _ any) (map[string]any, error) {
		return map[string]any{
			"databases.test-storage.driver": "db-driver-mock",
			"databases.test-storage.source": "...",
		}, nil
	})
	return factory.Build()
}
