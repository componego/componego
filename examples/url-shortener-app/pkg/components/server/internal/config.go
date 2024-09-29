package internal

import (
	"time"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/config"
	"github.com/componego/componego/impl/processors"
)

type Config struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	StopTimeout       time.Duration
}

func NewConfig(env componego.Environment) *Config {
	return &Config{
		Addr: config.GetOrPanic[string]("server.addr", processors.Multi(
			processors.DefaultValue(":3030"),
			processors.ToString(),
		), env),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		StopTimeout:       10 * time.Second,
	}
}
