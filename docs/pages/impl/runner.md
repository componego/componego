# Application Runner

## Basic Information

This entity is responsible for running the [application](./application.md) using the [driver](./driver.md).
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
This line in the main function is sufficient to start your [application](./application.md).

The function ^^runner.RunAndExit^^ runs the application and exits the program with an exit code after the application stops.

!!! note
    If the application completes with an error, the exit code will not equal 0 (componego.SuccessExitCode).

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

There are also methods ^^RunWithContext^^ and ^^RunGracefullyAndExit^^, which allow you to run the application using a custom context or facilitate a graceful shutdown.

## Application Mode

As you can see, you can run the application in different modes. By default, there are several modes available:

| Name                     | Description                |
|--------------------------|----------------------------|
| componego.ProductionMode | for production environment |
| componego.DeveloperMode  | for developers             |
| componego.TestMode       | for application tests      |

But you can add any mode you want.

You can retrieve the application mode through the [environment](./environment.md):
    ```go
    if env.ApplicationMode() == componego.DeveloperMode {
        // ...
    }
    ```

!!! note
    You should always use production mode when running your application on a production server.
    It is also recommended to use test mode when executing the application in [tests](../tests/runner.md#test-mode).

!!! note
    We strongly recommend using multiple application entry points, as demonstrated in [this example](https://github.com/componego/componego/tree/master/examples/url-shortener-app/cmd/application){:target="_blank"}.

    We believe that applications should be aware of the mode in which they will be launched even before execution.
    For example, this approach allows you to [read](./config.md#configuration-reader) different configurations based on the environment,
    rather than constructing the environment according to the configuration.

## Custom Runner

The custom runner is significant as it serves as an entry point where you can begin modifying the core of the framework to meet your specific requirements.

!!! note
    If you are a beginner with our framework, please skip this section
    and return to it after you have thoroughly read the rest of the documentation.

### Specific Driver Options

This is related to the application [driver](./driver.md), but you can manage it through the runner.
The options are factories that implement all default entities provided by the framework.
So this is the key (but not the only one) how you can replace the core of the framework with your code.

For example, you can pass additional options to your application.
Let's create a new Run function that accepts arguments:
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

    func Run(app componego.Application, appMode componego.ApplicationMode, additionalOptions any) int {
        d := driver.New(&driver.Options{
            Additional: additionalOptions,
            // ... other options
        })
        exitCode, err := d.RunApplication(context.Background(), app, appMode)
        if err != nil {
            _, _ = fmt.Fprint(os.Stderr, unhandled_errors.ToString(err, appMode, unhandled_errors.GetHandlers()))
        }
        return exitCode
    }
    ```

Review the options which are available in the driver code to ensure you can control everything.

!!! note
    The ability to replace the core of the framework is crucial because it allows for flexibility and customization.
    This means you are not bound to the default implementations of certain functions in the framework.
    Instead, you can substitute them with alternative methods that comply with the required interfaces.
    This feature enhances the adaptability of your application, enabling you to tailor functionalities to meet specific needs
    or to integrate with other systems more effectively.

    However, even more critical is the ability to replace business logic easily, as this functionality can be instrumental in creating mocks and other testing scenarios.
    The framework supports this capability as well. We've covered the details on [other pages](../tests/mock.md#rewriting-rules),
    highlighting how you can substitute specific business logic implementations to facilitate testing and improve code maintainability.

### Errors Handing

!!! note
    It is recommended to use a special application method [ApplicationErrorHandler](./application.md#applicationerrorhandler) to catch global errors or panic.

At the runner level you can handle errors that were not handled at all previous levels.

Based on the previous example, the following lines can be added for error handling:
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
You can convert the error into a string using specific handlers and display it to users.
This approach allows for a user-friendly presentation of error messages.
You can customize the error handling logic to fit your application’s needs.
Here's how you might implement this:
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

In addition to function ^^AddBefore^^, t, there are numerous other functions designed for handling errors in a specific sequence using an ordered map.
For a complete list of these methods, you can refer to the [source code here](https://github.com/componego/componego/tree/master/libs/ordered-map){:target="_blank"}.

!!! note
    Don't hesitate to explore the core of the framework and copy methods to modify them for your specific requirements.
    However, it's important to adhere to the [rewriting rules](../tests/mock.md#rewriting-rules) outlined by the framework.
