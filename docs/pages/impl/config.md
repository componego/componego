# Configuration

## Basic Information

Configuration is an important part of every application.

The framework has a single point for reading the configuration.
This is a special method that you can add to the [application](./application.md#applicationconfiginit) struct:
    ```go hl_lines="12 22"
    package application

    import (
        "github.com/componego/componego"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode) (map[string]any, error) {
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
As you can see, this method returns a map with configuration keys and values.

You can also return an error if there was an error reading the configuration.

!!! note
    Since the method accepts the [environment](./environment.md), you can return different configurations depending on the [mode](./runner.md#application-mode) in which the application is running.

This method is called only once and should return the configuration for the [application](./application.md) and [all components](./component.md) within that application.


## Configuration Reader

You can read the configuration in different ways as you like.

For example, you can use third-party libraries to get the configuration.
However, your function or library must return a variable of type ^^map[string]any^^.
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

    func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode) (map[string]any, error) {
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

You can also add post-processing of values after reading the configuration:
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
in:
    ```json
    {
        "server": {
            "addr": ":1234"
        }
    }
    ```
You can also use the default value after pipe:
    ```json
    {
        "server": {
            "addr": ":${ENV:APP_PORT|1234}"
        }
    }
    ```

We have described the map that function ^^ApplicationConfigInit^^ will return.
The next section describes how to get this value in your [application](./application.md) or [component](./component.md).

## Configuration Getter

The configuration value can be obtained using the [environment](./environment.md) in any part of the application:
    ```go
    value, err := env.ConfigProvider().ConfigValue("server.addr", nil)
    ```

Golang is a strongly typed language and it is impossible to use generics in this case in the current version of the language (go1.22).

There is an additional function inside the framework that helps solve the typing problem:
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
The shorter code looks like this:
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
    We use a dot as a separator between configuration keys for different levels of nesting.

## Configuration Processor

Validation and transformations of configuration values can be done using [processors](./processor.md).
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
    Typing is an important part of processors. The generic must meet the processor. Otherwise, there may be an error.
    ```go hl_lines="3"
    config.GetOrPanic[int64]("server.port", processors.Multi(
        processors.IsRequired(),
        processors.ToInt64(),
    ), env)
    ```
    We recommend that you always use a processor to change the type, because there is no guarantee that ^^ApplicationConfigInit^^ will return a value of the type you want.

## Configuration Examples

We recommend always creating an example configuration file when you create an [application](./application.md) or [component](./component.md).

For example, you created a component and a file with an example configuration of this component.

A developer who will use this component will be able to copy and merge the example configuration file into his main application configuration file.

After this, this file will be read in function ^^ApplicationConfigInit^^.
