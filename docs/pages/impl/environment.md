# Application Environment

## Basic information

An environment package serves as a comprehensive solution for managing the runtime of application environment.
This package not only provides access to [dependency injection (DI)](./dependency.md) management but also includes features for
handling application [active components](./component.md),
provides access to [configuration](./config.md), [application mode](./runner.md#application-mode),
[IO](./environment.md#application-io) and [global context](./environment.md#application-context).

## How to get environment

The environment can be accessed in several ways.
The first and easiest way is to get this object in the [application action](./application.md#applicationaction).
    ```go
    func (a *Application) ApplicationAction(env componego.Environment, _ []string) (int, error) {
        env.GetContext()
        env.Application()
        env.ApplicationIO()
        env.ApplicationMode()
        env.Components()
        env.DependencyInvoker()
        env.ConfigProvider()

        // ...
        return componego.SuccessExitCode, nil
    }

    var _ componego.Application = (*Application)(nil)
    ```

You can also get this object via DI.
    ```go
    type MyType struct {
        env componego.Environment `componego:"inject"`
    }

    func (m *MyType) Method() {
        if m.env.ApplicationMode() == componego.DeveloperMode {
            // ...
        }
    }
    ```
or
    ```go
    err := dependencyInvoker.Invoke(func(env componego.Environment) {
        // ...
    })
    ```

or another ways described in the [documentation about DI](./dependency.md).

!!! note
    The environment object cannot be rewritten by [rewritten DI objects](./dependency.md#rewriting-dependencies).
    This object is present in any application.


## How to use environment

| Method                  | Description                                                | Documentation                                  |
|-------------------------|------------------------------------------------------------|------------------------------------------------|
| env.GetContext()        | returns a current application context                      | [open](./environment.md#application-context)   |
| env.SetContext(newCtx)  | sets a new application context                             | [open](./environment.md#application-context)   |
| env.Application()       | returns a current application object                       | [open](./application.md)                       |
| env.ApplicationIO()     | returns an object for getting application input and output | [open](./environment.md#application-io)        |
| env.ApplicationMode()   | returns the mode in which the application is started       | [open](./runner.md#application-mode)           |
| env.Components()        | returns a sorted list of active application components     | [open](./component.md)                         |
| env.DependencyInvoker() | returns an object to invoke dependencies                   | [open](./dependency.md#access-to-dependencies) |
| env.ConfigProvider()    | returns an object for getting config                       | [open](./config.md#configuration-getter)       |


This is a universal key for accessing any part of the application.

## Application Context

It is recommended to use the application context to run various functions.
You can also replace the current context with another context, but the new context must inherit from the previous main context.

Let's create an example of a graceful shutdown [component](./component.md):
    ```go hl_lines="16 27"
    package component

    import (
        "context"
        "os"
        "os/signal"
        "syscall"

        "github.com/componego/componego"
        "github.com/componego/componego/impl/managers/component"
    )

    func NewComponent() componego.Component {
        factory := component.NewFactory("graceful-shutdown", "graceful-shutdown@example", "0.0.1")
        factory.SetComponentInit(func(env componego.Environment) error {
            ctx, cancel := context.WithCancel(env.GetContext())
            interruptChan := make(chan os.Signal, 1)
            signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
            go func() {
                select {
                case <-interruptChan:
                case <-ctx.Done():
                    signal.Stop(interruptChan)
                }
                cancel()
            }()
            return env.SetContext(ctx)
        })
        return factory.Build()
    }
    ```

Components functions and how to create them are described on the [next documentation page](./component.md).

The framework is also provide a [graceful shutdown component](../qcomponents/graceful-shutdown.md), created according to the same example.

## Application IO

If you want to output (or receive) some text to (from) the console, then you must use special methods:
    ```go hl_lines="2"
    func (a *Application) ApplicationAction(env componego.Environment, _ []string) (int, error) {
        appIO := env.ApplicationIO()
        _, _ = fmt.Fprintln(appIO.OutputWriter(), "your text")
        _, _ = fmt.Fprintln(appIO.ErrorOutputWriter(), "your error text")
        reader := bufio.NewReader(appIO.InputReader())
        text, _ := reader.ReadString('\n')
    }
    ```
