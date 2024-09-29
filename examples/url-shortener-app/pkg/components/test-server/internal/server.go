package internal

import (
	"net/http"
	"net/http/httptest"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency"

	"github.com/componego/componego/examples/url-shortener-app/third_party/servermux"
)

type TestServer struct {
	env componego.Environment
}

func NewServer(env componego.Environment) *TestServer {
	return &TestServer{
		env: env,
	}
}

func (s *TestServer) Run(
	configureCallback func(addRouter func(pattern string, handler any)),
	runCallback func(baseUrl string),
) {
	mainRouter, addRouter := servermux.CreateRouter()
	configureCallback(func(pattern string, handler any) {
		castedHandler, err := dependency.Invoke[http.HandlerFunc](handler, s.env)
		if err != nil {
			panic(err)
		}
		addRouter(pattern, castedHandler)
	})
	serverInstance := httptest.NewServer(mainRouter)
	defer serverInstance.Close()
	runCallback(serverInstance.URL)
}
