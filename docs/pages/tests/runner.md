# Tests Runner

## Basic Example

This is the simplest way:
    ```go hl_lines="6 12-13"
    package tests

    import (
        "testing"

        "github.com/componego/componego/tests/runner"

        "secret.com/project-x/tests/mocks"
    )

    func TestExample(t *testing.T) {
        env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock(), nil)
        t.Cleanup(cancelEnv)
        // ... here you can use application environment.
    }
    ```
In the example above, we created the new [environment](../impl/environment.md) based on the [application mock](./mock.md).
You can use this environment to run the necessary functions in your tests.

The last argument of the function accepts the test options, including [driver options](../impl/runner.md#specific-driver-options), that will be applied to the current test.

When the environment is canceled, all necessary functions will be called to stop the application.

The framework is thread-safe so you can run tests in parallel.
However, your personal code or the code of third-party libraries you use may not be thread-safe.

!!! note
    In the example above, you can see how to test an application.

    To test [component](../impl/component.md), you must create an application that [depends](../impl/application.md#applicationcomponents) only on that component.
    For example, this way you can configure a component, because [this method](../impl/config.md) is in any application.

## Test Mode

In tests, the application should be launched in ^^componego.TestMode^^.
There are different [application launch modes](../impl/runner.md#application-mode).
However, for tests, it is recommended to use ^^componego.TestMode^^.

The code above runs the application in this mode, so please consider this detail in your implementation.
