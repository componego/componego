package test_server

import (
	"github.com/componego/componego"

	"github.com/componego/componego/examples/url-shortener-app/pkg/components/server"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/test-server/internal"
)

type (
	TestServer = internal.TestServer
)

type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

// ComponentIdentifier belongs to interface componego.Component.
func (c *Component) ComponentIdentifier() string {
	return "componego:examples:test-server"
}

// ComponentVersion belongs to interface componego.Component.
func (c *Component) ComponentVersion() string {
	return "0.0.1"
}

// ComponentComponents belongs to interface componego.ComponentComponents.
func (c *Component) ComponentComponents() ([]componego.Component, error) {
	return []componego.Component{
		server.NewComponent(),
	}, nil
}

// ComponentDependencies belongs to interface componego.ComponentDependencies.
func (c *Component) ComponentDependencies() ([]componego.Dependency, error) {
	return []componego.Dependency{
		internal.NewServer,
	}, nil
}

var (
	_ componego.Component             = (*Component)(nil)
	_ componego.ComponentComponents   = (*Component)(nil)
	_ componego.ComponentDependencies = (*Component)(nil)
)
