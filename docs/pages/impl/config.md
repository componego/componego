# Configuration

## Basic Information

Configuration is an important part of every application.

The framework provides a single point for reading the configuration through a special method that you can add to an [application](./application.md#applicationconfiginit) struct:
    ```go hl_lines="12 22"
    package application

    import (
        "github.com/componego/componego"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode, options any) (map[string]any, error) {
        return map[string]any{
            "config.key1": "config.value1",
            "config.key2": "config.value2",
            // ...
        }, nil
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationConfigInit = (*Application)(nil)
    )
    ```
This method returns a map containing the configuration keys and values.
You can also return an error if there was an issue reading the configuration.

This method is called only once and should return the configuration for the [application](./application.md) and [all components](./component.md) within that application.

## Configuration Reader

You can read the configuration in various ways.

For example, you can use third-party libraries to obtain the configuration.
However, your function or library must return a variable of the type ^^map[string]any^^:
    ```go hl_lines="7 18 25"
    package application

    import (
        "fmt"

        "github.com/componego/componego"
        "github.com/spf13/viper"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode, options any) (map[string]any, error) {
        switch appMode {
        case componego.ProductionMode:
            return ConfigReader("./config/production.config.json")
        // ...
        default:
            return nil, fmt.Errorf("not supported application mode: %d", appMode)
        }
    }

    func ConfigReader(filename string) (map[string]any, error) {
        // This function should be placed in a separate package.
        v := viper.New()
        v.SetConfigFile(filename)
        if err := v.ReadInConfig(); err != nil {
            return nil, err
        }
        return v.AllSettings(), nil
    }

    var (
        _ componego.Application           = (*Application)(nil)
        _ componego.ApplicationConfigInit = (*Application)(nil)
    )
    ```

You can also perform post-processing of values after reading the configuration:
    ```go hl_lines="5 13"
    package config

    import (
        "github.com/componego/componego/impl/environment/managers/config"
        "github.com/spf13/viper"
    )

    func ConfigReader(filename string) (map[string]any, error) {
        // This function should be placed in a separate package.
        v := viper.NewWithOptions()
        // ...
        settings := v.AllSettings()
        return settings, config.ProcessVariables(settings)
    }

    ```

This function converts the following values:
    ```json
    {
        "server": {
            "addr": ":${ENV:APP_PORT}"
        }
    }
    ```
into:
    ```json
    {
        "server": {
            "addr": ":1234"
        }
    }
    ```
You can also use the default value after a pipe:
    ```json
    {
        "server": {
            "addr": ":${ENV:APP_PORT|1234}"
        }
    }
    ```

We have described the map that ^^ApplicationConfigInit^^ returns.
The next section describes how to get this value in your [application](./application.md) or [component](./component.md).

## Configuration Getter

Any configuration values can be accessed using the [environment](./environment.md) in any part of the application:
    ```go
    value, err := env.ConfigProvider().ConfigValue("server.addr", nil)
    ```

Golang is a strongly typed language and it is impossible to use generics ^^in this case^^ in the current version of the language (go1.22).

There is an additional function within the framework that helps solve the typing problem:
    ```go hl_lines="5 9"
    package config

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/config"
    )

    func GetServerAddr(env componego.Environment) (string, error) {
        return config.Get[string]("server.addr", nil, env)
    }
    ```
The shorter code appears as follows:
    ```go hl_lines="5 9"
    package config

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/config"
    )

    func GetServerAddr(env componego.Environment) string {
        return config.GetOrPanic[string]("server.addr", nil, env)
    }
    ```

!!! note
    We use a dot as a separator between configuration keys to indicate different levels of nesting.

## Configuration Processor

Validation and transformation of configuration values can be performed using [processors](./processor.md):
    ```go hl_lines="6"
    package config

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/config"
        "github.com/componego/componego/impl/processors"
    )

    func GetServerAddr(env componego.Environment) string {
        return config.GetOrPanic[string]("server.addr", processors.Multi(
            processors.DefaultValue(":3030"),
            processors.ToString(),
        ), env)
    }
    ```

!!! note
    Typing is a crucial aspect of processors. The generic type must match the processor. Otherwise, there may be an error:
    ```go hl_lines="3"
    config.GetOrPanic[int64]("server.port", processors.Multi(
        processors.IsRequired(),
        processors.ToInt64(),
    ), env)
    ```
    You should use a processor to change the type, as there is no guarantee that ^^ApplicationConfigInit^^ will return a value of the desired type.

## Configuration Struct

Application configurations can indeed become quite large, and managing each configuration key with individual [processors](./processor.md) can be inefficient.
Instead, it's recommended to define a struct that contains fields corresponding to different configuration values.
You can validate this struct in a single step during the initialization of the dependency injection [DI container](./dependency.md).

Each component can have its own separate struct to describe its configuration.
This allows for better modularity and separation of concerns.
You can find an example of such a struct-based configuration approach [here](https://medium.com/@konstanchuk/25bfd16a97a9#413f){:target="_blank"}.

## Configuration Examples

It’s recommended to create an example configuration file when you create an [application](./application.md) or [component](./component.md).

For instance, if you create a component, you should also include a file with an example of its configuration.
This allows developers who use the component to easily copy and merge the example configuration file into their main application configuration file.
After that, it can be read in the ^^ApplicationConfigInit^^ function.
