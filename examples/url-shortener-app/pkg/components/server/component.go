package server

import (
	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/pkg/components/server/internal"
)

type (
	Server = internal.Server
)

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

// ComponentIdentifier belongs to interface componego.Component.
func (c *Component) ComponentIdentifier() string {
	return "componego:examples:server"
}

// ComponentVersion belongs to interface componego.Component.
func (c *Component) ComponentVersion() string {
	return "0.0.1"
}

// ComponentDependencies belongs to interface componego.ComponentDependencies.
func (c *Component) ComponentDependencies() ([]componego.Dependency, error) {
	return []componego.Dependency{
		internal.NewServer,
		internal.NewConfig,
	}, nil
}

var (
	_ componego.Component             = (*Component)(nil)
	_ componego.ComponentDependencies = (*Component)(nil)
)
