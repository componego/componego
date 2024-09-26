# Application Mock

## Basic Example

If you want to replace a part of your [application](../impl/application.md) for testing, you can do so easily.

Use the example of the application below to see how to create a mock:
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
We use inheritance. In a child struct (mock), you can override methods to return new values.

Here's an example of overriding dependencies, but you can add any other method:
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
We call the parent method and then add new data.

It is explained how to run the mock on the [next documentation page](./runner.md).

## Rewriting Rules

We've already described this on previous documentation pages for each framework entity, but let's sum it up.

!!! note
    We follow the [order in which the methods are called](../impl/driver.md#application-initialization-order).
    The last methods called may override the return data of the previous methods.

### For Components

If several [components](../impl/component.md#componentidentifier) the same identifier, the component specified last will take precedence.

### For Dependencies

The last [dependencies](../impl/dependency.md) in the list are applied.

If you use a [constructor](../impl/dependency.md#dependency-constructors) for dependencies, the last constructor will overwrite the dependencies if its return types match those of the constructor you want to replace.

You will encounter an error if you return a constructor that produces a different set of types, except for the last returned type if it is an error.

If you use an object directly instead of a constructor, you can overwrite this value with either the same object or a constructor that returns an object of the same type.

!!! note
    Objects returned by the [environment](../impl/environment.md#how-to-use-environment) methods cannot be overwritten..

### For Configuration

Everything is straightforward here. You need to read data from another resource.

For example, you can use the [application mode](../impl/runner.md#application-mode) to read various configurations:
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

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode, options any) (map[string]any, error) {
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

You can also override this method in a mock to return a map with a different configuration required for your test.

### For Framework Core

If you want to replace something within the framework, you can use [driver options](../impl/runner.md#specific-driver-options).

Don't hesitate to copy core functions to make the necessary changes for your project.
We would also appreciate any [any suggestion or pull request](../contribution/guide.md) that can help develop this project.
