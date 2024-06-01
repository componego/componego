package config

import (
	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/config"
	"github.com/componego/componego/impl/processors"
)

func GetServerAddr(env componego.Environment) string {
	return config.GetOrPanic[string]("server.addr", processors.Multi(
		processors.DefaultValue(":3030"),
		processors.ToString(),
	), env)
}
