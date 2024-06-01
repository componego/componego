# Tests Runner

## Basic Example

The simplest way is this:
    ```go hl_lines="6 12-13"
    package tests

    import (
        "testing"

        "github.com/componego/componego/tests/runner"

        "secret.com/project-x/tests/mocks"
    )

    func TestExample(t *testing.T) {
        env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock())
        t.Cleanup(cancelEnv)
        // ... here you can use application environment.
    }
    ```
In this example, we created a new [environment](../impl/environment.md) based on the [application mock](./mock.md).
You can use this environment to run the necessary functions in your tests.

When the environment is canceled, all necessary functions will be called to stop the application.

The framework is thread safe so you can run tests in parallel.
However, your personal code or the code of third-party libraries you use may not be thread safe.

!!! note
    The example above showed how to test an application.

    To test [component](../impl/component.md), you must create an application that [depends](../impl/application.md#applicationcomponents) only on that component.
    This way you can, for example, configure a component, because [this method](../impl/config.md) is in the application.

## Test Mode

In tests, the application should be launched in mode ^^componego.TestMode^^.
There are different [application launch modes](../impl/runner.md#application-mode).
However, for tests, it is recommended to use mode ^^componego.TestMode^^.

The above code runs the application in this mode. Take this subtlety into account in your code.
