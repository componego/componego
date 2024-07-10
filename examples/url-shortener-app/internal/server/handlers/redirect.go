package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/componego/componego/examples/url-shortener-app/internal/repository"
	"github.com/componego/componego/examples/url-shortener-app/internal/server/json"
	"github.com/componego/componego/examples/url-shortener-app/internal/server/validation"
)

// -------------------- GET Redirect -------------------

func NewRedirectGetHandler(redirectRepository repository.RedirectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if redirect, err := redirectRepository.Get(r.Context(), key); err == nil {
			http.Redirect(w, r, redirect.Url, http.StatusPermanentRedirect)
			return
		}
		http.NotFound(w, r)
	}
}

// -------------------- PUT Redirect -------------------

func NewRedirectPutHandler(redirectRepository repository.RedirectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Url string `json:"url"`
		}

		type response struct {
			NewUrl string `json:"newUrl"`
		}

		requestData, err := json.Get[request](r)
		if err != nil {
			json.Send(w, nil, errors.New("error decoding JSON"))
			return
		}
		if !validation.IsValidUrl(requestData.Url) {
			json.Send(w, nil, errors.New("invalid url"))
			return
		}

		if redirect, err := redirectRepository.Add(r.Context(), requestData.Url); err != nil {
			json.Send(w, nil, errors.New("failed to add redirect"))
		} else if r.TLS == nil {
			// noinspection ALL
			json.Send(w, &response{
				NewUrl: fmt.Sprintf("http://%s/get/%s", r.Host, redirect.Key),
			}, nil)
		} else {
			json.Send(w, &response{
				NewUrl: fmt.Sprintf("https://%s/get/%s", r.Host, redirect.Key),
			}, nil)
		}
	}
}
