package server

import (
	"errors"

	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/internal/server/handlers"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/server"
)

func Run(env componego.Environment, s *server.Server) error {
	if err := errors.Join(
		s.AddRouter("GET /", handlers.NewIndexGetHandler),
		s.AddRouter("PUT /create", handlers.NewRedirectPutHandler),
		s.AddRouter("GET /get/{key}", handlers.NewRedirectGetHandler),
	); err != nil {
		return err
	}
	return s.Run(env.GetContext())
}
