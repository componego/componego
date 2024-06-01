package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/componego/componego/examples/url-shortener-app/internal/domain"
	"github.com/componego/componego/examples/url-shortener-app/internal/utils"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database"
)

type RedirectRepository interface {
	Add(ctx context.Context, url string) (*domain.Redirect, error)
	Get(ctx context.Context, key string) (*domain.Redirect, error)
}

type redirectRepository struct {
	db *sql.DB
}

func NewRedirectRepository(dbProvider database.Provider) (RedirectRepository, error) {
	db, err := dbProvider.Get("main-storage")
	if err != nil {
		return nil, err
	}
	return &redirectRepository{
		db: db,
	}, nil
}

func (r *redirectRepository) Add(ctx context.Context, url string) (*domain.Redirect, error) {
	randomString := utils.GetRandomString(10)
	if randomString == "" {
		return nil, errors.New("could not get the new url key")
	}
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO redirects(url_key, url) VALUES(?, ?)`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	redirect := &domain.Redirect{
		Key: randomString,
		Url: url,
	}
	if _, err = stmt.ExecContext(ctx, redirect.Key, redirect.Url); err != nil {
		return nil, err
	}
	return redirect, nil
}

func (r *redirectRepository) Get(ctx context.Context, key string) (*domain.Redirect, error) {
	stmt, err := r.db.PrepareContext(ctx, `SELECT url FROM redirects WHERE url_key = ?`)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	redirect := &domain.Redirect{
		Key: key,
	}
	if err = stmt.QueryRowContext(ctx, key).Scan(&redirect.Url); err != nil {
		return nil, err
	}
	return redirect, nil
}
