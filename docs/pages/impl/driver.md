# Application Driver

## Basic information

The application driver plays a main role in Componego runtime by serving as the entry point that initializes
and orchestrates all essential functions within an application.
Essentially, it acts as the engine that kick-starts the execution of an application by coordinating various components,
managing configurations, and initiating critical processes.
This driver ensures a smooth and controlled stage setting for seamless operation.

## Differences between Runner

The difference is minimal, but the driver can be shared among many applications,
initiating the basic functions of the application based on the driver's options.
Driver options control various aspects of the application, including the environment factory, dependency manager,
configuration manager, input reader, output writer, and error output writer.
These options are flexible and can be modified since they are only [options](./runner.md#specific-driver-options).

In most cases, you don't need to be aware of the driver options. However, if you wish to modify any core aspects of the framework,
you can explore the [source code](https://github.com/componego/componego/tree/master/impl/driver){:target="_blank"} to see how it is implemented.

## Application initialization order

<figure markdown>
  ![Componego Flow](../assets/images/diagrams/componego-flow.svg){ width="622" height="1202" loading=lazy }
  <figcaption>Componego Flow</figcaption>
</figure>

!!! note
    The red elements in the image can handle errors that occur in previous (or nested) functions.

!!! note
    We recommend looking at this diagram again when you fully understand
    how to create [applications](./application.md) and [components](./application.md), and the entities they provide.

The general order in which functions are called is as follows:

1. runner.Run
2. driver.RunApplication
3. application.ApplicationConfigInit
4. application.ApplicationComponents
5. component.ComponentComponents (+ getting components for each component)
6. component.ComponentDependencies (for each of the active components)
7. application.ApplicationDependencies
8. component.ComponentInit (for each of the active components)
9. application.ApplicationAction
10. component.ComponentStop (for each of the active components in reverse order)
11. application.ApplicationErrorHandler (If there was an error)
12. exit

Not all methods are described here (if the [application](./application.md) or [component](./component.md) uses these methods).
This list provides a sufficient overview of the application initialization order.

!!! note
    The order of initialization and method calls is crucial when rewriting elements of the application.
    For example, an [application](./application.md)  can rewrite [dependencies](./dependency.md) of [component](./component.md)
    because a method that returns dependencies for the application object (^^ApplicationDependencies^^) is called
    after the same function for components (^^ComponentDependencies^^).
    This behavior can be particularly useful when creating [mocks](../tests/mock.md).
