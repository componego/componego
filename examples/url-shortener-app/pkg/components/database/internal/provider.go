package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/componego/componego"
)

type Provider interface {
	GetConnection(name string) (*sql.DB, error)
	CreateConnection(name string) (*sql.DB, error)
	CloseConnection(name string) error
}

type provider struct {
	mutex sync.Mutex
	env   componego.Environment
	list  map[string]*sql.DB
}

func NewProvider(env componego.Environment) Provider {
	dbProvider := &provider{
		mutex: sync.Mutex{},
		env:   env,
		list:  make(map[string]*sql.DB, 2),
	}
	return dbProvider
}

func (p *provider) GetConnection(name string) (*sql.DB, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.list == nil {
		// The application is stopping now.
		// The function to close all provider was called.
		return nil, fmt.Errorf("you cannot create a connection to '%s'. Make sure the order of the components is correct", name)
	} else if connection, ok := p.list[name]; ok {
		return connection, nil
	} else if connection, err := p.CreateConnection(name); err == nil {
		p.list[name] = connection
		return connection, nil
	} else { //nolint:revive
		return nil, err
	}
}

func (p *provider) CreateConnection(name string) (db *sql.DB, err error) {
	var driver, source string
	if driver, err = getDriver(name, p.env); err != nil {
		return nil, err
	}
	if source, err = getDataSourceName(name, p.env); err != nil {
		return nil, err
	}
	db, err = sql.Open(driver, source)
	if err != nil {
		return nil, err
	} else if err = db.PingContext(p.env.GetContext()); err != nil {
		return nil, err
	}
	return db, nil
}

func (p *provider) CloseConnection(name string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if connection, ok := p.list[name]; ok {
		delete(p.list, name)
		return connection.Close()
	}
	return fmt.Errorf("not found connection with name '%s'", name)
}

func (p *provider) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	errs := make([]error, 0, len(p.list))
	for _, connection := range p.list {
		errs = append(errs, connection.Close())
	}
	// It sets a flag that the connection can no longer be opened.
	p.list = nil
	// We switch the runtime so that waiting goroutines can complete database provider.
	runtime.Gosched()
	return errors.Join(errs...)
}

var (
	_ Provider  = (*provider)(nil)
	_ io.Closer = (*provider)(nil)
)
