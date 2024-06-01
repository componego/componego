package internal

import (
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/config"
	"github.com/componego/componego/impl/processors"
)

func getDataSourceName(connectionName string, env componego.Environment) (string, error) {
	return config.Get[string](fmt.Sprintf("databases.%s.source", connectionName), processors.Multi(
		processors.IsRequired(),
		processors.ToString(),
	), env)
}

func getDriver(connectionName string, env componego.Environment) (string, error) {
	return config.Get[string](fmt.Sprintf("databases.%s.driver", connectionName), processors.Multi(
		processors.IsRequired(),
		processors.ToString(),
	), env)
}
