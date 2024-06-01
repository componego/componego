package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency"

	"github.com/componego/componego/examples/url-shortener-app/third_party/errgroup"
	"github.com/componego/componego/examples/url-shortener-app/third_party/servermux"

	"github.com/componego/componego/examples/url-shortener-app/internal/config"
	"github.com/componego/componego/examples/url-shortener-app/internal/server/handlers"
)

type route struct {
	pattern string
	handler any
}

func CreateRouter(env componego.Environment) (http.HandlerFunc, error) {
	return injectRouters(env, []route{
		// GoLang >= v1.22
		{"GET /", handlers.NewIndexGetHandler},
		{"PUT /create", handlers.NewRedirectPutHandler},
		{"GET /get/{key}", handlers.NewRedirectGetHandler},
	})
}

func Run(env componego.Environment) error {
	router, err := CreateRouter(env)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              config.GetServerAddr(env),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           router,
	}
	return listenAndServe(env, server)
}

func injectRouters(env componego.Environment, routers []route) (http.HandlerFunc, error) {
	router, add := servermux.CreateRouter()
	for _, route := range routers {
		// All the handlers will have the desired dependencies without having to pass them directly.
		handler, err := dependency.Invoke[http.HandlerFunc](route.handler, env)
		if err != nil {
			return nil, err
		}
		add(route.pattern, handler)
	}
	return router, nil
}

func listenAndServe(env componego.Environment, server *http.Server) error {
	writer := env.ApplicationIO().OutputWriter()
	_, err := fmt.Fprintf(writer, "server is starting on %s...\n", server.Addr)
	if err != nil {
		return err
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
		// We use a graceful shutdown component so the server will be stopped gracefully if it is using the environment context.
		case <-env.GetContext().Done():
		case <-stopChan:
			return nil
		}
		_, err := fmt.Fprintf(writer, "trying to stop the server on %s...\n", server.Addr)
		cancelableCtx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()
		return errors.Join(server.Shutdown(cancelableCtx), err)
	})
	return errGroup.Wait()
}
