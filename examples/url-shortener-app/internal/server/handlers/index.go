package handlers

import (
	"fmt"
	"net/http"
)

func NewIndexGetHandler() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(response, "It works")
	}
}
