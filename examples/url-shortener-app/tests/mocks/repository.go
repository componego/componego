package mocks

import (
	"context"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/environment/managers/dependency"

	"github.com/componego/componego/examples/url-shortener-app/internal/domain"
	"github.com/componego/componego/examples/url-shortener-app/internal/repository"
)

type redirectRepositoryMock struct {
	originalRepository repository.RedirectRepository
}

func NewRedirectRepositoryMock(env componego.Environment) (repository.RedirectRepository, error) {
	originalRepository, err := dependency.Invoke[repository.RedirectRepository](repository.NewRedirectRepository, env)
	if err != nil {
		return nil, err
	}
	return &redirectRepositoryMock{
		originalRepository: originalRepository,
	}, nil
}

func (r *redirectRepositoryMock) Add(ctx context.Context, url string) (*domain.Redirect, error) {
	// In this case, this mock looks like a proxy. This is just an example.
	// You can create a real mock based on this example.
	return r.originalRepository.Add(ctx, url)
}

func (r *redirectRepositoryMock) Get(ctx context.Context, key string) (*domain.Redirect, error) {
	return r.originalRepository.Get(ctx, key)
}
