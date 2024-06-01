# Graceful-Shutdown Component

The framework provides a component for smoothly stopping the application functions.

Connect this component like any other component as follows:

=== "In Application"
    ```go hl_lines="5 13 22"
    package application

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/qcomponents/graceful-shutdown"
    )

    type Application struct {
    }

    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{
            graceful_shutdown.NewComponent(),
            // ...
        }, nil
    }

    // ...

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationComponents = (*Application)(nil)
    )
    ```
=== "In Component"
    ```go hl_lines="5 13 22"
    package component

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/qcomponents/graceful-shutdown"
    )

    type Component struct {
    }

    func (c *Component) ComponentComponents() ([]componego.Component, error) {
        return []componego.Component{
            graceful_shutdown.NewComponent(),
            // ...
        }, nil
    }

    // ...

    var (
        _ componego.Component           = (*Component)(nil)
        _ componego.ComponentComponents = (*Component)(nil)
    )
    ```

Now you can use a [context](../impl/environment.md#application-context) that will be canceled if a stop signal is received:
    ```go hl_lines="2"
    go func() {
		<-env.GetContext().Done()
		service.Shutdown()
	}()
	service.Run()
    ```

It should be taken into account that functionality of graceful-shutdown applies only to the next [components](../impl/component.md) in the list
and [application](../impl/application.md). For example:
    ```go hl_lines="5"
    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{
            component1.NewComponent(), // without graceful-shutdown context.
            component2.NewComponent(), // without graceful-shutdown context.
            graceful_shutdown.NewComponent(),
            component3.NewComponent(), // with graceful-shutdown context.
            component4.NewComponent(), // with graceful-shutdown context.
            // ...
        }, nil
    }
    ```
