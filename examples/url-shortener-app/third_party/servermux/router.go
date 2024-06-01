package servermux

import (
	"net/http"
)

type AddHandler func(pattern string, handler func(http.ResponseWriter, *http.Request))

type router struct {
	base *http.ServeMux
}

func CreateRouter() (http.HandlerFunc, AddHandler) {
	r := &router{
		base: http.NewServeMux(),
	}
	return r.ServeHTTP, r.base.HandleFunc
}

func (r *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	r.base.ServeHTTP(response, request)
}
