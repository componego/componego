# Processors

## Basic information

Processors are special functions that convert variables of one type into another and validate them.

On the [previous page](./config.md#configuration-processor) you have seen the example of using processors:
    ```go hl_lines="6 12"
    package config

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/config"
        "github.com/componego/componego/impl/processors"
    )

    func GetConnectionName(env componego.Environment) string {
        return config.GetOrPanic[string](
            "connection.name",
            processors.IsRequired(),
            env,
        )
    }
    ```

## Default Processors

| Function                          | Description                                         |
|-----------------------------------|-----------------------------------------------------|
| processors.ToBool()               | converts the value to boolean                       |
| processors.IsBool()               | checks whether the value is a boolean value         |
| processors.ToInt64()              | converts the value to int64                         |
| processors.ToFloat64()            | converts the value to float64                       |
| processors.ToString()             | converts the value to string                        |
| processors.IsRequired()           | checks that the value is present (not nil)          |
| processors.DefaultValue(anyValue) | sets the default value if the previous value is nil |

## Custom Processor

There are 2 simple ways to create it:
=== "Long code"
    ```go hl_lines="15"
    package processor

    import (
        "github.com/componego/componego"
    )

    type Processor struct{}

    func (p *Processor) ProcessData(value any) (any, error) {
        // convert a value to another or validate a value
        // if there is an error, then you must return it
        return newValue, nil
    }

    var _ componego.Processor = (*Processor)(nil)
    ```
=== "Short code"
    ```go hl_lines="5"
     package processor

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/processors"
    )

    func NewProcessor() componego.Processor {
        return processors.New(func(value any) (any, error) {
            // convert a value to another or validate a value
            // if there is an error, then we must return it
            return newValue, nil
        })
    }
    ```
!!! note
    Like in other parts of the application, you can use [dependencies](./dependency.md) within processors:
    ```go
    package processor

    import (
        "github.com/componego/componego"
    )

    type Processor struct {
        env componego.Environment `componego:"inject"`
    }

    // ...
    ```

## Multi Processor

This can be implemented using ^^processors.Multi^^, which allows for combining multiple functions into one:
    ```go hl_lines="11"
    package processor

    import (
        "strings"

        "github.com/componego/componego"
        "github.com/componego/componego/impl/processors"
    )

    func NewProcessor() componego.Processor {
        return processors.Multi(
            processors.ToString(),
            processors.New(func(value any) (any, error) {
                return strings.Split(value.(string), ","), nil
            }),
        )
    }
    ```
Look at the example above: we do not check if the value is a ^^string^^ in the second processor,
as we are confident it will be converted to a ^^string^^ in the first processor.

This is a chain of function calls executed sequentially.
If an error occurs in any function, the chain will be interrupted.

## Processor As Validator

You can use processors not only to convert data to another format but also for validation:
    ```go
    package processor

    import (
        "fmt"

        "github.com/componego/componego"
        "github.com/componego/componego/impl/processors"
    )

    func ToAge() componego.Processor {
        return processors.Multi(
            processors.ToInt64(),
            processors.New(func(value any) (any, error) {
                age := value.(int64)
                if age < 21 {
                    return nil, fmt.Errorf("invalid age: %d", age)
                }
                return age, nil
            }),
        )
    }
    ```
You can also use third-party libraries to validate specific rules:
    ```go hl_lines="15"
    package processor

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/processors"
        "github.com/go-playground/validator/v10" // third-party library
    )

    var validate = validator.New()

    func ToAge() componego.Processor {
        return processors.Multi(
            processors.ToInt64(),
            processors.New(func(value any) (any, error) {
                return value, validate.Var(value, "required,numeric,min=21")
            }),
        )
    }
    ```

!!! note
    The framework is designed as a way to run an application.
    You should avoid using framework functions within the main loop of your application.
    Instead, utilize the framework solely for initializing the main application loop.
