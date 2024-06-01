---
social_meta:
    title: Get Started
---
The Componego Framework seamlessly integrates component, dependency injection, configuration, event and error handling,
providing a robust foundation for building modular, scalable, and maintainable software systems.

At its core, the framework embraces a component-based architecture.
These components encapsulate specific functionalities and seamlessly integrate to construct larger, more complex applications.
Promoting code reusability and maintainability, the component framework allows developers to work on independent units,
fostering efficient development and a cohesive application architecture.

Main features of the framework:

1. well-organized application initialization process and smallest main function
2. simple and powerful components
3. flexible dependency injections without code generation
4. the framework core does not depend on third-party packages
5. easy integration with previously written code and third-party packages
6. no impact on the business logic code of your application. You will write and organize your business logic as you did before
7. any entity of the framework can be replaced without changing previously written code
8. good documentation for developers

The framework does not provide a database connection, web server, queue and other things.
Instead, we provide the ability to conveniently integrate these things into components for later reuse in your projects.
You can use absolutely anything you want. There are no limits on packages.
Any existing Golang package can be easily wrapped in our components.

!!! note
    The documentation is a brief description of the main functions of the framework.
    You may find some things complicated, but for each part that is described in another section, we add links.

    You need to create an application or component using this framework to fully understand the framework.
    Don't be afraid to open the [source code](https://github.com/componego/componego){:target="_blank"}. There you will find many interesting functions that were not described in the documentation.

<section class="get-started-menu" markdown="block">
<div class="componego-grid-cards" markdown="block">
<div class="card" markdown="block">
:octicons-star-fill-16:[How to create an application?](./impl/application.md)
</div>
<div class="card" markdown="block">
:octicons-play-16:[How to run an application?](./impl/runner.md)
</div>
<div class="card" markdown="block">
:octicons-database-16:[How to create a component?](./impl/component.md)
</div>
<div class="card" markdown="block">
:octicons-git-pull-request-16:[How to use dependency injections?](./impl/dependency.md)
</div>
<div class="card" markdown="block">
:octicons-plug-16:[How to use configuration?](./impl/config.md)
</div>
<div class="card" markdown="block">
:octicons-bug-16:[How to handle errors?](./impl/application.md#applicationerrorhandler)
</div>
<div class="card" markdown="block">
:octicons-copy-16:[How to create mocks?](./tests/mock.md)
</div>
<div class="card" markdown="block">
:octicons-terminal-16:[How to create tests?](./tests/runner.md)
</div>
</div>
</section>
