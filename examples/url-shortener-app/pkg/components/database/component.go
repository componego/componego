package database

import (
	"errors"

	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database/internal"
)

type (
	Provider = internal.Provider
)

type Component struct {
	closeConnections func() error
}

func NewComponent() *Component {
	return &Component{}
}

// ComponentIdentifier belongs to interface componego.Component.
func (c *Component) ComponentIdentifier() string {
	return "componego:examples:database"
}

// ComponentVersion belongs to interface componego.Component.
func (c *Component) ComponentVersion() string {
	return "0.0.1"
}

// ComponentDependencies belongs to interface componego.ComponentDependencies.
func (c *Component) ComponentDependencies() ([]componego.Dependency, error) {
	return []componego.Dependency{
		func(env componego.Environment) Provider {
			dbProvider, closeConnections := internal.NewProvider(env)
			c.closeConnections = closeConnections
			return dbProvider
		},
	}, nil
}

func (c *Component) ComponentStop(_ componego.Environment, prevErr error) error {
	// We check that the function exists, since the function may not exist if the dependency has been replaced.
	if c.closeConnections != nil {
		return errors.Join(prevErr, c.closeConnections())
	}
	return prevErr
}

var (
	_ componego.Component             = (*Component)(nil)
	_ componego.ComponentDependencies = (*Component)(nil)
	_ componego.ComponentStop         = (*Component)(nil)
)
