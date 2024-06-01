# Application Runner

## Basic Information

This is the entity that runs the [application](./application.md) using the [driver](./driver.md).
    ```go hl_lines="11"
    package main

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/runner"

        "github.com/componego/componego/examples/hello-app/internal/application"
    )

    func main() {
        runner.RunAndExit(application.New(), componego.ProductionMode)
    }
    ```
This line in the main function is enough to start your [application](./application.md).

Function ^^runner.RunAndExit^^ starts the application and after stopping the application exits the program with exit code.

!!! note
    If the application completed with an error, then the exit code from the application will not be equal to 0 (componego.SuccessExitCode).

You can also use ^^runner.Run^^, which starts the application but does not exit it:
    ```go hl_lines="11"
    package main

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/runner"

        "github.com/componego/componego/examples/hello-app/internal/application"
    )

    func main() {
        exitCode := runner.Run(application.New(), componego.ProductionMode)
        // ...
    }
    ```

## Application Mode

As you can see, you can run the application in different modes.

By default, there are several modes:

| Name                     | Description                |
|--------------------------|----------------------------|
| componego.ProductionMode | for production environment |
| componego.DeveloperMode  | for developers             |
| componego.TestMode       | for application tests      |

But you can add any mode you want.

You can get the application mode through the [environment](./environment.md).
    ```go
    if env.ApplicationMode() == componego.DeveloperMode {
        // ...
    }
    ```

!!! note
    You should always use production mode when you run your application on a production server.
    It is also recommended to use test mode when you run the application in [tests](../tests/runner.md#test-mode).

!!! note
    We strongly recommend using multiple application entry points,
    as shown in this [example](https://github.com/componego/componego/tree/master/examples/url-shortener-app/cmd/application){:target="_blank"}.

    We believe that applications should know in what mode it will be launched even before launching.
    For example, this allows you to [read](./config.md#configuration-reader) different configurations depending on the environment,
    instead of building the environment depending on the configuration.

## Custom Runner

The custom runner is significant because it is where you can start modifying the framework core to suit your specific requirements.

!!! note
    If you are new to our framework, please skip this section.
    Come back to it after you've fully read the rest of the documentation pages.

### Specific Driver Options

This is related to the application [driver](./driver.md), but you can control it through the runner.
Options are some factories that implement all the entities that the framework provides.
So this is the key (but not the only one) how you can replace the core of the framework with your code.

Fo example, the runner uses the Golang application arguments to run (^^os.Args^^). You can specify your custom ones.
Let's create a new Run function that takes arguments.
    ```go hl_lines="15"
    package custom_runner

    import (
        "context"
        "fmt"
        "os"

        "github.com/componego/componego"
        "github.com/componego/componego/impl/driver"
        "github.com/componego/componego/impl/runner/unhandled-errors"
    )

    func Run(app componego.Application, appMode componego.ApplicationMode, args []string) int {
        d := driver.New(&driver.Options{
            Args: args,
            // ... other options
        })
        exitCode, err := d.RunApplication(context.Background(), app, appMode)
        if err != nil {
            _, _ = fmt.Fprint(os.Stderr, unhandled_errors.ToString(err, appMode, unhandled_errors.GetHandlers()))
        }
        return exitCode
    }
    ```

!!! note
    We are creating a context in this code.
    You can read about how to use this context [here](./environment.md#application-context).

!!! note
    The ability to replace the core of the framework is important because you are not tied to the implementation of some functions of the framework.
    You can replace them with other methods that satisfy some interfaces.

    However, more important is the ability to easily replace business logic, because this can be used in mocks and more.
    The framework can do this too. We have described this on [other pages](../tests/mock.md#rewriting-rules).

### Errors Handing

!!! note
    It is recommended to use a special application] method [ApplicationErrorHandler](./application.md#applicationerrorhandler) to catch global errors or panic.

At the runner level you can handle errors that were not handled at all previous levels.

Based on the previous example, the following lines add error handling.
    ```go hl_lines="5-6"
    func Run(app componego.Application, appMode componego.ApplicationMode) int {
        d := driver.New(nil)
        exitCode, err := d.RunApplication(context.Background(), app, appMode)
        if err != nil {
            handlers := unhandled_errors.GetHandlers()
            _, _ = fmt.Fprint(os.Stderr, unhandled_errors.ToString(err, appMode, handlers))
        }
        return exitCode
    }
    ```
We convert the error into a string using some handlers and show it to the user.
Of course, you can use any error handling you want or add your handlers to the standard handlers.
    ```go hl_lines="2"
    handlers := unhandled_errors.GetHandlers()
    handlers.AddBefore(
        "company:my-handler-name",
        func(err error, writer io.Writer, appMode componego.ApplicationMode) bool {
            if errors.Is(err, MyError) {
                // ...
                return true // return true if the error is processed
            }
            return false
        },
        "componego:vendor-proxy",
    )
    _, _ = fmt.Fprint(os.Stderr, unhandled_errors.ToString(err, appMode, handlers))
    ```

In addition to function ^^AddBefore^^, there are many other functions for handling errors in the required sequence (ordered map).
Look at the [source code](https://github.com/componego/componego/tree/master/libs/ordered-map){:target="_blank"} for a complete list of these methods.

!!! note
    Don't be afraid to look into the core of the framework and copy methods to make changes for your specific requirements.
    However, try to follow the [rewriting rules](../tests/mock.md#rewriting-rules) provided by the framework.
