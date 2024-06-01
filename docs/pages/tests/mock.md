# Application Mock

## Basic Example

Let's say you want to replace some part of your [application](../impl/application.md) for testing.
This is very easy to do and is one of the main features of this framework.

Let's use an example application to look at how you can create a mock.
    ```go hl_lines="10 15 22"
    package mocks

    import (
        "github.com/componego/componego"

        "secret.com/project-x/internal/application"
    )

    type ApplicationMock struct {
        *application.Application
    }

    func NewApplicationMock() *ApplicationMock {
        return &ApplicationMock{
            Application: application.New(),
        }
    }

    // ... other methods

    var (
        _ componego.Application = (*ApplicationMock)(nil)
        // ...
    )
    ```
As you can see, we use inheritance. In a child (mock) struct, you can rewrite methods and return new values.

Here's an example of what it might look like for [dependencies](./../impl/dependency.md), but you can add any other method.
    ```go  hl_lines="2 3 10"
    func (a *ApplicationMock) ApplicationDependencies() ([]componego.Dependency, error) {
        dependencies, err := a.Application.ApplicationDependencies()
        dependencies = append(
            dependencies,
            func() Service {
                return &mockService{}
            },
            // ...
            )
        return dependencies, err
    }
    ```
We call the parent method and add new data.

How to run the mock is shown on the [next documentation page](./runner.md).

## Rewriting Rules

We have already described this on previous documentation pages for each framework entity. But let's sum it up.

!!! note
    We use the [order in which the methods are called](../impl/driver.md#application-initialization-order).
    Methods that were called last may rewrite the return data of previous methods.

### For Components

If several [components](../impl/component.md#componentidentifier) have the same identifier, the component specified last will be applied.

### For Dependencies

The last [dependencies](../impl/dependency.md) in the list are applied.

If you use a [constructor](../impl/dependency.md#dependency-constructors) for dependencies, the dependencies will be rewritten by the last constructor if the return types match the constructor you want to rewrite.

You will get an error if you return a constructor that returns a different set of types. The exception is the last returned type if it is an error.

If instead of a constructor you use an object directly, then you can rewrite this value with the same object or constructor that returns the same object type.

!!! note
    Objects that are returned by [environment](../impl/environment.md#how-to-use-environment) methods cannot be rewritten.

### For Configuration

Everything is simple here. You need to read data from another resource.

For example, you can use [application mode](../impl/runner.md#application-mode) to read different configurations:
    ```go hl_lines="18 20 22"
    package application

    import (
        "fmt"

        "github.com/componego/componego"
        // ...
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode) (map[string]any, error) {
        switch appMode {
        case componego.ProductionMode:
            return config_reader.Read("./config/production.config.json")
        case componego.DeveloperMode:
            return config_reader.Read("./config/developer.config.json")
        case componego.TestMode:
            return config_reader.Read("./config/test.config.json")
        default:
            return nil, fmt.Errorf("not supported application mode: %d", appMode)
        }
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationConfigInit = (*Application)(nil)
    )
    ```

You can also rewrite this method in a mock and return a map with a different configuration that is needed for your test.

### For Framework Core

If you want to replace something inside the framework, then you can use [driver options](../impl/runner.md#specific-driver-options).

Don't be afraid to copy core functions to make the changes necessary for your project.
We will also be glad to receive [any suggestion or pull request](../contribution/guide.md) that will help develop and improve this project.
