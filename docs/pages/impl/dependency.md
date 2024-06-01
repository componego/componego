# Dependency Injection

## Basic Information

Dependency injection is a design pattern used in software engineering to achieve inversion of control between classes and their dependencies.
In simpler terms, it's a technique where the dependencies of a class are provided from the outside rather than created within the class itself.
This helps to decouple components and promotes easier testing, maintainability, and flexibility in your code.

Dependencies can be provided by [components](./component.md#componentdependencies) and the [application](./application.md#applicationdependencies). Special methods inside these entities are responsible for this.
Let's look at an example:

=== "In Application"
    ```go hl_lines="14-16 23"
    package application

    import (
        "github.com/componego/componego"
    )

    type Application struct {
    }

    // ...

    func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{
            func() SomeService {
                return &someService{}
            },
            // ...
        }, nil
    }

    var (
        _ componego.Application             = (*Application)(nil)
        _ componego.ApplicationDependencies = (*Application)(nil)
    )
    ```
=== "In Component"
    ```go hl_lines="14-16 23"
    package component

    import (
        "github.com/componego/componego"
    )

    type Component struct {
    }

    // ...

    func (c *Component) ComponentDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{
            func() SomeService {
                return &someService{}
            },
            // ...
        }, nil
    }

    var (
        _ componego.Component             = (*Component)(nil)
        _ componego.ComponentDependencies = (*Component)(nil)
    )
    ```
Now you can use this provided object in your application.

It is recommended to use constructors to create dependencies. Consider these constructors further on this page.

## Dependency Constructors

Let's look at the following code example and possible variations of what the constructor might look like:
    ```go hl_lines="3"
    func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{
            NewProductRepository,
            // ...
        }, nil
    }

    // ...
    ```

1. The constructor returns a struct as a pointer:
        ```go
        func NewProductRepository() *ProductRepository {
            return &ProductRepository{}
        }
        ```
2. The constructor returns a struct as an interface:
        ```go
        func NewProductRepository() ProductRepository {
            return &productRepository{}
        }
        ```
3. The constructor can return an error as the latest value:
       ```go
       func NewProductRepository() (ProductRepository, error) {
           return &productRepository{}, nil
       }
       ```
4. The constructor can accept an unlimited number of dependencies:
       ```go
       func NewProductRepository(db * database.Provider) ProductRepository {
           return &productRepository{
               db: db,
           }
       }
       ```
5. You can even do things like this:
       ```go
       func NewProductRepository(di componego.DependencyInvoker) (ProductRepository, error) {
           repo := &productRepository{}
           return repo, di.PopulateFields(repo)
       }
       ```
!!! note
    Constructors can accept and return an unlimited number of dependencies. However, they must be presented as pointers.

    It is also recommended to use interfaces, as it can be convenient in some cases.

    The constructor, like any framework entity, is thread-safe.

Another alternative is to represent the dependency directly as an object:
    ```go hl_lines="3"
    func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{
            &ProductRepository{},
            // ...
        }, nil
    }

    // ...
    ```
Most likely you will rarely use this method.

!!! note
    Loops between dependencies are not allowed. You will receive an error message when starting the application if a loop occurs.

## Access to Dependencies

Dependencies can be obtained in several ways. The easiest way is to use the [environment](./environment.md).

### Invoke

This is a method that takes a function as an argument, which can take any dependencies
provided in the [components](./component.md#componentdependencies) or [application](./application.md#applicationdependencies).
    ```go
    _, err := env.DependencyInvoker().Invoke(func(service SomeService, repository SomeRepository) {
        // ...
    })
    ```
The function may also return an error as the last return value:
    ```go
    _, err := env.DependencyInvoker().Invoke(func(service SomeService) error {
        // ...
        return service.Action()
    })
    ```
The called function can return a value:
    ```go
    returnValue, err := env.DependencyInvoker().Invoke(func(service SomeService) int {
        // ...
        return service.Action()
    })
    // or
    returnValue, err := env.DependencyInvoker().Invoke(func(service SomeService) (int, error) {
        // ...
        return service.Action()
    })
    ```
!!! note
    Since the return type is type ^^any^^, you can use a helper to get the correct type:
    ```go hl_lines="5 9"
    package example

    import (
        "github.com/componego/componego"
        "github.com/componego/componego/impl/environment/managers/dependency"
    )

    func GetValue(env componego.Environment) (int, error) {
        intValue, err := dependency.Invoke[int](SomeFunction, env)
        // intValue := dependency.InvokeOrPanic[int](SomeFunction, env)
        return intValue, err
    }
    ```

You can also get an object for dependency injection inside a function:
    ```go
    _, err := di.Invoke(func(di componego.DependencyInvoker, service SomeService) (any, error) {
        // ...
        return di.Invoke(service.Action)
    })
    ```
Although in this case, you could use closures.

### Populate

This is a function that fills a variable, which is a pointer.
    ```go
    var service *Service
    err := env.DependencyInvoker().Populate(service)
    ```
!!! note
    The type of the variable must exactly match the requested type.

    The difference between functions ^^Populate^^ and ^^Invoke^^ is that the first function can only accept a struct because only a struct can be a pointer.
    At the same time, the second function can accept arguments of any type included in the list of allowed types for dependencies.


### PopulateFields

The ^^Populate^^ method fills only a variable, but much more often you need to fill fields in a struct.
Let's look at an example:
    ```go hl_lines="2"
    type Service struct {
        dbProvider database.Provider `componego:"inject"`
    }

    // ...

    service := &Service{}
    err := env.DependencyInvoker().PopulateFields(service)
    ```
This method fills only those fields that have the special tag shown in the example. All other fields are ignored.
Fields can be private or public. The field type can be any.

If an error occurs, the method will return it.

## Default Dependencies

Each application has a set of standard dependencies with which you can access various functions of the application.

The table shows a set of these dependencies:

| Variable                        | Description                                                              |
|---------------------------------|--------------------------------------------------------------------------|
| env componego.Environment       | access to the [environment](./environment.md)                            |
| app componego.Application       | returns the current [application](./application.md)                      |
| appIO componego.ApplicationIO   | access to the [application IO](./environment.md#application-io)          |
| di componego.DependencyInvoker  | returns the [dependency invoker](./dependency.md#access-to-dependencies) |
| config componego.ConfigProvider | provides access to [configuration](./config.md#configuration-getter)     |

These are objects that are returned by the [environment](./environment.md) using its methods.

!!! note
    Although you can get [context](./environment.md#application-context) through the environment, you cannot get context through dependencies.
    Use the environment directly to obtain the application context.

!!! note
    Standard dependencies cannot be rewritten. You must use the [driver options](./runner.md#specific-driver-options) if you want to change them.

## Rewriting Dependencies

Rewriting is one of the main features of the framework. Let's see how you can rewrite dependencies.
    ```go hl_lines="4 7"
    func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
        return []componego.Dependency{
            func() SomeService {
                return &someService{}
            },
            func() SomeService {
                return &someService2{}
            },
            // ...
        }, nil
    }

    // ...
    ```
In this case, the second service will be used because it is specified after the constructor of the first service.

For rewriting to work, the return value types must match.  This is the only condition.
Constructors can accept any dependency, but the return types must match. Only in this case rewriting will work.

If you try to return the type that was returned in the constructors above, but the returned type does not match, you will receive an error.

!!! note
    The only exception is the last type returned if that type is an error.

Rewriting dependencies is one of the key elements in creating [mocks](../tests/mock.md) using this framework.

Remember that according to the documentation about the [order of initialization of elements](./driver.md#application-initialization-order) of the framework,
method ^^ApplicationDependencies^^ is called after the same function for [components](./component.md#componentdependencies) (^^ComponentDependencies^^).

This means that you can rewrite dependencies in your [application](./application.md#applicationdependencies) that were declared in [components](./component.md#componentdependencies).
You can also rewrite dependencies in components that were added in [parent components](./component.md#componentcomponents).
