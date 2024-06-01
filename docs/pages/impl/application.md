# Application

## Basic Application

The application is a struct with a set of methods that describe the elements of the application.

Let's try to create your first application and launch it.
    ```go
    package main

    import (
        "fmt"

        "github.com/componego/componego"
        "github.com/componego/componego/impl/application"
        "github.com/componego/componego/impl/runner"
    )

    type Application struct {
    }

    func (a *Application) ApplicationName() string {
        return "Hello World App"
    }

    func (a *Application) ApplicationAction(env componego.Environment, args []string) (int, error) {
        _, err := fmt.Fprintln(env.ApplicationIO().OutputWriter(), "Hello World!")
        return application.ExitWrapper(err)
    }

    func main() {
        runner.RunAndExit(&Application{}, componego.ProductionMode)
    }
    ```
Save this code to a file and run it without arguments.

Don't forget to initialize [go mod](https://go.dev/ref/mod#go-mod-init){:target="_blank"} and [download all dependencies](https://go.dev/ref/mod#go-mod-tidy){:target="_blank"}.
    ```text hl_lines="1"
    % go run main.go
    Hello world!
    ```

The source code for this application example is available [here](https://github.com/componego/componego/tree/master/examples/hello-app){:target="_blank"}.

About runner is written in the [next page](./runner.md).

## Mandatory methods

### ApplicationName

The function returns the application name.
    ```go
    func (a *Application) ApplicationName() string {
        return "The Best Application"
    }

    // ...
    ```

### ApplicationAction

The function describes the main action of the current application.
    ```go
    func (a *Application) ApplicationAction(env componego.Environment, args []string) (int, error) {
        // ...
        return componego.SuccessExitCode, nil
    }

    // ...
    ```
In this function you can write business logic for your application.

## Optional methods

### ApplicationComponents

The application may depend on [components](./component.md):
    ```go
    func (a *Application) ApplicationComponents() ([]componego.Component, error) {
        return []componego.Component{ /* ... */ }, nil
    }

    // ...
    ```

### ApplicationDependencies

It can provide [dependencies](./dependency.md):
    ```go
    func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{ /* ... */ }, nil
    }

    // ...
    ```

### ApplicationConfigInit

The application can read the [configuration](./config.md):
    ```go
    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode) (map[string]any, error) {
        return map[string]any{
            "config.key": "config.value",
        }, nil
    }

    // ...
    ```
You can return different configuration depending on the [mode](./runner.md#application-mode) the application is running in.

### ApplicationErrorHandler

By default, there is also a method that handles all your errors that were not handled previously:
    ```go
    func (a *Application) ApplicationErrorHandler(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error {
        if errors.Is(err, MyError) {
            // ...
            err = nil
        } else {
            // ...
        }
        return err
    }
    ```
This method also intercepts global panic in the application.

Unhandled errors returned by this method will be received and processed by the [runner](./runner.md#errors-handing) at the core level.

You can also catch errors at the [component level](./component.md#componentstop).

<hr/>

!!! note
    For more clarity and compile time validation you can add the following code:
    ```go
    var (
        _ componego.Application             = (*Application)(nil)
        _ componego.ApplicationComponents   = (*Application)(nil)
        _ componego.ApplicationDependencies = (*Application)(nil)
        _ componego.ApplicationConfigInit   = (*Application)(nil)
        _ componego.ApplicationErrorHandler = (*Application)(nil)
    )
    ```
    The names of interface correspond to the logic they implement.

    It is recommended to always add such validation in order to easily find and fix problems in the code
    in case changes are made to interface methods in future versions of the framework.

!!! note
    You can also add your own methods and implement your custom logic to them, because an application is a struct that just implements interfaces.

!!! note
    Pay special attention to the [order in which methods are called](./driver.md#application-initialization-order).
    This will help you understand the application initialization process.

## Application Factory

There is also a short code for creating the application.
However, we do not recommend using this method of creating an application.
    ```go
    package main

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/application"
        "github.com/componego/componego/impl/runner"
    )

    func main() {
        factory := application.NewFactory("Application Name")
        factory.SetApplicationDependencies(func() ([]componego.Dependency, error) {
            return []componego.Dependency{ /* ... */ }, nil
        })
        // ... other methods.
        runner.RunAndExit(factory.Build(), componego.ProductionMode)
    }
    ```
The factory has different Set* methods that correspond to the methods described above.

## Application Skeleton

You can quickly create an application skeleton using the following ways:
    ```shell
    curl -sSL https://raw.githubusercontent.com/componego/componego/master/tools/create-basic-app.sh | sh
    ```
or
    ```shell
    wget -O - https://raw.githubusercontent.com/componego/componego/master/tools/create-basic-app.sh | sh
    ```

On Windows, you can run the above commands with [Git Bash](https://git-scm.com/download/win){:target="_blank"}.

This command will create a ^^componego-basic-app^^ folder with the most basic version of the application, based on which you can begin development.

An example of a full-fledged application using our framework can be found [here](https://github.com/componego/componego/tree/master/examples/url-shortener-app){:target="_blank"}.

To learn more, visit other documentation pages.
