---
social_meta:
    title: Get Started
---
Componego Framework seamlessly integrates components, dependency injection, configuration, and error handling,
providing a robust foundation for building modular, scalable, and maintainable software systems.

The framework embraces a component-based architecture.
These components encapsulate specific functionalities and seamlessly integrate to construct more complex applications.
By promoting code reusability and maintainability, the framework allows developers to work on independent units.

Main features of the framework:

1. A well-organized application initialization process with a minimal main function.
2. Simple and powerful components.
3. Flexible dependency injection that doesn't require code generation.
4. The core of the framework does not depend on third-party packages.
5. Easily integrates with existing code and third-party packages.
6. No impact on the business logic of your application. You can write and organize it just as you did before.
7. Any entity of the framework can be replaced without changing previously written code.
8. Comprehensive and clear documentation for developers.

The framework does not provide a database connection, web server, queue, and other features.
Instead, we offer the ability to conveniently integrate these components for later reuse in your projects.
You can use any packages you prefer; there are no limits.
Any existing Golang package can be easily wrapped in the framework's components.

!!! note
    The documentation provides a brief description of the main functions of the framework.
    You may find some aspects complicated, but we include links for each part described in another section.

    You need to create an application or component using this framework to fully understand it.
    Don't hesitate to open the [source code](https://github.com/componego/componego){:target="_blank"}.
    There, you will find many interesting functions that are not described in the documentation.

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
