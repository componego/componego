package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency"

	"github.com/componego/componego/examples/url-shortener-app/third_party/errgroup"
	"github.com/componego/componego/examples/url-shortener-app/third_party/servermux"
)

type Server struct {
	env       componego.Environment
	config    *Config
	router    http.Handler
	addRouter servermux.AddHandler
}

func NewServer(env componego.Environment, config *Config) *Server {
	router, addRouter := servermux.CreateRouter()
	return &Server{
		env:       env,
		config:    config,
		router:    router,
		addRouter: addRouter,
	}
}

func (s *Server) AddRouter(pattern string, handler any) error {
	castedHandler, err := dependency.Invoke[http.HandlerFunc](handler, s.env)
	if err != nil {
		return err
	}
	s.addRouter(pattern, castedHandler)
	return nil
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.debug("server is starting on %s...\n", s.config.Addr); err != nil {
		return err
	}
	server := &http.Server{
		Addr:              s.config.Addr,
		ReadTimeout:       s.config.ReadTimeout,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		Handler:           s.router,
	}
	stopChan := make(chan struct{}, 1)
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		defer func() {
			stopChan <- struct{}{}
		}()
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	errGroup.Go(func() error {
		select {
		case <-ctx.Done():
		case <-stopChan:
			return nil
		}
		err := s.debug("trying to stop the server on %s...\n", s.config.Addr)
		cancelableCtx, cancelCtx := context.WithTimeout(context.Background(), s.config.StopTimeout)
		defer cancelCtx()
		return errors.Join(server.Shutdown(cancelableCtx), err)
	})
	return errGroup.Wait()
}

func (s *Server) debug(format string, args ...any) error {
	_, err := fmt.Fprintf(s.env.ApplicationIO().OutputWriter(), format, args...)
	return err
}
