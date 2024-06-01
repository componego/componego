# Application Components

## Basic Component

A component is a modular, reusable piece of code that performs a specific function in a larger application.
The components are designed as independent, self-contained units that can be easily integrated into an [application](./application.md).

Components are not microservices. These are folders that contain different functionality.
To describe what functionality a component depends on and what functionality the component provides, we use a struct with a set of different methods.

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
This component doesn't do anything because it lacks the methods the various functions provide.

Look at the methods that the component can provide in the documentation below in the next section.
For now, let's look at how you can add a component to your [application](./application.md#applicationcomponents):
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

This function returns the component ID.
    ```go
    func (c *Component) ComponentIdentifier() string {
        return "company-name:component-name"
    }

    // ...
    ```
If the identifier in several components is the same, then the last component will be used.
This can be used to [rewrite components](./component.md#initialization-rewriting).

### ComponentVersion

This function returns the component version.
    ```go
    func (c *Component) ComponentVersion() string {
        return "0.0.1"
    }

    // ...
    ```
The version must match the [SemVer](https://semver.org/){:target="_blank"} format.

!!! note
    It is recommended to always increase the component version if you make changes to the component.

## Optional methods

!!! note
    The component struct is similar to the [application](./application.md) struct.
    However, there are slight differences in the methods.

### ComponentComponents

The component may depend on components:
    ```go
    func (a *Component) ComponentComponents() ([]componego.Component, error) {
        return []componego.Component{ /* ... */ }, nil
    }

    // ...
    ```

!!! note
    The order in which the components are loaded depends on the dependencies between the components.

    This also affects the rewriting of components.
    The components that are listed last rewrite the components that were added first (if they have the same identifier).

### ComponentDependencies

Like an [application](./application.md#applicationdependencies), it can provide [dependencies](./dependency.md):
    ```go
    func (a *Component) ComponentDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{ /* ... */ }, nil
    }

    // ...
    ```

### ComponentInit

This is a function that is called in order of dependencies between components.
You can add some of your own logic to the method.
    ```go
    func (a *Component) ComponentInit(env componego.Environment) error {
        // ...
        return nil
    }

    // ...
    ```

### ComponentStop

This method is called before component stop.
You can handle previous error (return a new or old error).
    ```go
    func (a *Component) ComponentStop(env componego.Environment, prevErr error)  error {
        // ...
        return prevErr
    }

    // ...
    ```

<hr/>

!!! note
    For more clarity and compile time validation you can add the following code:
    ```go
    var (
        _ componego.Component             = (*Component)(nil)
        _ componego.ComponentComponents   = (*Component)(nil)
        _ componego.ComponentDependencies = (*Component)(nil)
        _ componego.ComponentInit         = (*Component)(nil)
        _ componego.ComponentStop         = (*Component)(nil)
    )
    ```
    The names of interface correspond to the logic they implement.

    It is recommended to always add such validation in order to easily find and fix problems in the code
    in case changes are made to interface methods in future versions of the framework.

## Initialization & Rewriting

Take a look at the following example:
    ```go hl_lines="17 19"
    package application

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/qcomponents/database"

        myDatabase "secret.com/project-x/components/database"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{
            database.NewComponent(),
            // ...
            myDatabase.NewComponent(),
        }, nil
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationComponents = (*Application)(nil)
    )
    ```
If both connected components have the same identifier, then only the last component will be used.

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
    As you can see, we print the unique ID in the loop so you can see which components are currently in use.

Let's look at a more complex example where components can depend on other components.

For a better understanding, let's assume that each component color represents a component with the same identifier.
The same component can be reused by other components at different levels of [dependencies between them](./component.md#componentcomponents).

<figure markdown>
  ![Components Init First Example](../assets/images/diagrams/components-init-flow.svg)
  <figcaption>First Example of Components Initialization</figcaption>
</figure>

We see that this tree of components has been converted into such a sorted list.
This is a list of active components in our application.

You can see here how components are rewritten.

!!! note
    In fact, the sorted list can be anything. This depends on the implementation of the component manager
    (as we know from other documentation pages, absolutely everything in our framework can be replaced).
    However, any manager implementation must ensure that components are loaded based on the dependencies between the components.

    The framework guarantees that the order of loading components in your application will always be the same
    unless there are changes in the [dependencies between components](./component.md#componentcomponents).

Let's look at another example where we just moved the red component down.

<figure markdown>
  ![Components Init Second Example](../assets/images/diagrams/components-init-flow-2.svg)
  <figcaption>Second Example of Components Initialization</figcaption>
</figure>

As you can see, the list has not changed, but the components from other levels are active.
This is how component rewriting works.

In short:

1. Components can be rewritten if they have the [same identifier](./component.md#componentidentifier).
   Active components are components that are last in the list (taking into account the levels between components).
2. initialization order of components is always the same unless you change the list order of [child components](./component.md#componentcomponents).

## Component Factory

There is also a short code for creating the component.

This is the same thing as for the [application](./application.md#application-factory).
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
