# Security Policy

## Supported Versions

We are releasing security vulnerability fixes for the package's latest version(s).

The framework core doesn't depend on third-party packages and uses only the packages provided in the standard [GoLang](https://go.dev/) packages.
We monitor language security updates and try to fix vulnerabilities as quickly as possible.

The documentation uses [MkDocs](https://github.com/mkdocs/mkdocs) / [MkDocs Material](https://github.com/squidfunk/mkdocs-material) with additional Python plugins.
It does not affect the runtime of the application written using our framework.
We are not responsible for supporting McDocs and are currently using the free version of this software.
You should not run the documentation on your PC without using docker for more security.

Vulnerabilities in other parts of the framework will be fixed as soon as they are reported.

## Reporting a Vulnerability

Please report (suspected) security vulnerabilities to [t.me/konstanchuk](https://t.me/konstanchuk).

Please include as much of the information listed below as you can to help us better understand and resolve the issue:
* The type of issue.
* Full paths of source file(s) related to the manifestation of the issue.
* The location of the affected source code (tag/branch/commit or direct repository URL).
* Any special configuration is required to reproduce the issue.
* Step-by-step instructions to reproduce the issue.
* Proof-of-concept or exploit code (if possible).
* Impact of the issue, including how an attacker might exploit the issue.

This information will help us triage your report more quickly.

You will receive a response from us within 48 hours.
If the issue is confirmed, we will release a patch as soon as possible depending on complexity but historically within a few days.
