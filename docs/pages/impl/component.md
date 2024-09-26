# Application Components

## Basic Component

A component is a modular, reusable piece of code that performs a specific function in a larger application.
Components are designed as independent, self-contained units that can be easily integrated into an [application](./application.md).

Components are not microservices. They are folders that contain different functionalities.
We use a struct with a set of methods to describe the functionalities a component depends on and the functionalities it provides.

The basic component looks like this:
    ```go
    package component

    import (
        "github.com/componego/componego"
    )

    type Component struct {
    }

    func NewComponent() *Component {
        return &Component{}
    }

    func (c *Component) ComponentIdentifier() string {
        return "company-name:component-name"
    }

    func (c *Component) ComponentVersion() string {
        return "0.0.1"
    }

    var (
        _ componego.Component = (*Component)(nil)
    )
    ```
This component does not perform any actions because it lacks the methods that provide various functionalities.

Pay attention to methods that the component can provide in the documentation below in the next section.
Look at the example below to see how you can add a component to your [application](./application.md#applicationcomponents):
    ```go hl_lines="14 20"
    package application

    import (
        "github.com/componego/componego"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{
            NewComponent(),
        }, nil
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationComponents = (*Application)(nil)
    )
    ```

## Mandatory methods

### ComponentIdentifier

This function returns the component ID:
    ```go
    func (c *Component) ComponentIdentifier() string {
        return "company-name:component-name"
    }

    // ...
    ```
If the identifier in multiple components is the same, the last one will be used.
You can use it for [overwriting components](./component.md#initialization-rewriting).

### ComponentVersion

This function returns the version of the component:
    ```go
    func (c *Component) ComponentVersion() string {
        return "0.0.1"
    }

    // ...
    ```
The version must match the [SemVer](https://semver.org/){:target="_blank"} format.

!!! note
    It is always recommended to increment the component version when making changes to the component.

## Optional methods

!!! note
    Any component struct is similar to the [application](./application.md) struct.
    However, there are slight differences in their methods.

### ComponentComponents

Any component may depend on other components:
    ```go
    func (a *Component) ComponentComponents() ([]componego.Component, error) {
        return []componego.Component{ /* ... */ }, nil
    }

    // ...
    ```

!!! note
    The order in which components are loaded depends on their dependencies.

    This also affects the component overwriting rules:
    components listed last will overwrite those listed first if they share the same identifier.

### ComponentDependencies

Like an [application](./application.md#applicationdependencies), it can provide [dependencies](./dependency.md):
    ```go
    func (a *Component) ComponentDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{ /* ... */ }, nil
    }

    // ...
    ```

### ComponentInit

This function is called in the order of dependencies between components.
You can add custom logic to this method:
    ```go
    func (a *Component) ComponentInit(env componego.Environment) error {
        // ...
        return nil
    }

    // ...
    ```

### ComponentStop

This method is called when the component stops.
You can handle the previous error by returning either a new or the original error:
    ```go
    func (a *Component) ComponentStop(env componego.Environment, prevErr error)  error {
        // ...
        return prevErr
    }

    // ...
    ```

<hr/>

!!! note
    For greater clarity and compile-time validation, you can add the following code:
    ```go
    var (
        _ componego.Component             = (*Component)(nil)
        _ componego.ComponentComponents   = (*Component)(nil)
        _ componego.ComponentDependencies = (*Component)(nil)
        _ componego.ComponentInit         = (*Component)(nil)
        _ componego.ComponentStop         = (*Component)(nil)
    )
    ```
    These names correspond to the logic they implement.
    It is always recommended to add such validation to easily find and fix problems in the code,
    especially if changes are made to interface methods in future versions of the framework.

## Initialization & Rewriting

Pay attention to the following example:
    ```go hl_lines="17 19"
    package application

    import (
        "github.com/componego/componego"

        "secret.com/project-x/components/database1"
        "secret.com/project-x/components/database2"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{
            database1.NewComponent(),
            // ...
            database2.NewComponent(),
        }, nil
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationComponents = (*Application)(nil)
    )
    ```
If both connected components share the same identifier, only the last one will be used.

!!! note
    You can get a list of active components through the [environment](./environment.md).
    ```go
    activeAndSortedComponents := env.Components()
    for _, component := range activeAndSortedComponents {
        _, _ = fmt.Fprintf(
            env.ApplicationIO().OutputWriter(),
            "identifier: %s, version: %s, struct type: %T\n",
            component.ComponentIdentifier(), component.ComponentVersion(), component,
        )
    }
    ```
    We print the unique ID in the loop so you can see which components are currently in use.

It also supports overriding components that have multiple levels of nesting.
You can observe how it works in the [project code](https://github.com/componego/componego/blob/master/impl/environment/managers/component/tests/manager.go){:target="_blank"}.

## Component Factory

There is also a brief code snippet for creating the component.
This is the same thing as for an [application](./application.md#application-factory).
    ```go
    package component

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/component"
    )

    func NewComponent() componego.Component {
        factory := component.NewFactory("identifier", "0.0.1")
        factory.SetComponentInit(func(env componego.Environment) error {
            // ...
            return nil
        })
        // ... other methods.
        return factory.Build()
    }
    ```
