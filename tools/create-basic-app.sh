#!/usr/bin/env sh

# Copyright 2024 Volodymyr Konstanchuk and the Componego Framework contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset

PACKAGE_NAME="github.com/componego/componego"
PROJECT_NAME="componego-basic-app"

if ! command -v "go" >/dev/null; then
    echo "ERROR: please install GoLang and try again."
    exit 1
fi

if [ ! -p /dev/stdin ] && command -v "realpath" >/dev/null; then
    cd "$(dirname "$(realpath -- "$0")")"
fi

create=false
for arg in "$@"; do
    if [ "$arg" = "--create" ]; then
        create=true
    fi
done

directory=$(pwd)
while [ "$directory" != "/" ]; do
    if [ -e "$directory/go.mod" ]; then
        echo "ERROR: could not create a project in the current directory because a project is already created: $directory/go.mod"
        exit 1
    fi
    directory=$(dirname "$directory")
done
directory=$(pwd)

if [ -e "${PROJECT_NAME}" ]; then
    echo "ERROR: could not create the required files for the project because the directory or file already exists: $(pwd)}/${PROJECT_NAME}"
    exit 1
fi

if [ "$create" = false ]; then
    temp_dir=$(mktemp -d)
    cd "${temp_dir}"

    cat >"temp_app.go" <<EOF
package main

import (
    "github.com/componego/componego"
)

func main() {
    _ = componego.ProductionMode
}
EOF
    go mod init github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME
    go mod tidy
    go mod verify

    go_path=$(go env GOPATH)
    version=$(go list -m -f '{{.Version}}' "${PACKAGE_NAME}")

    temp_dir=$(mktemp -d)
    cd "${temp_dir}"

    cp "${go_path}/pkg/mod/${PACKAGE_NAME}@${version}/tools/create-basic-app.sh" "${temp_dir}/create-basic-app.sh"

    if [ "$version" = "v0.0.1" ]; then
        head -n210 "${temp_dir}/create-basic-app.sh" > "${temp_dir}/create-basic-app-temp.sh"
        cp -f "${temp_dir}/create-basic-app-temp.sh" "${temp_dir}/create-basic-app.sh"
    fi

    sh create-basic-app.sh --create

    cd "${PROJECT_NAME}"

    go mod init github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME
    go get "${PACKAGE_NAME}@${version}"
    go mod tidy
    go mod verify

    go fmt ./...

    # CGO_ENABLED=1 go test -race -v -count=1 ./...
    go test -v -count=1 ./... # without -race because CGO may not be installed.

    cp -r "${temp_dir}/${PROJECT_NAME}" "${directory}/${PROJECT_NAME}"
    echo "The project has been created. Output directory -> ${directory}/${PROJECT_NAME}"
    exit "$?"
fi

mkdir -p "${PROJECT_NAME}"
cd "${PROJECT_NAME}"

mkdir -p "cmd/application"
cat >cmd/application/main.go <<EOF
package main

import (
    "github.com/componego/componego"
    "github.com/componego/componego/impl/runner"

    "github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/application"
)

func main() {
    // This is an entry point for launching the application in production mode.
    runner.RunAndExit(application.New(), componego.ProductionMode)
}
EOF

mkdir -p "cmd/application/dev"
cat >cmd/application/dev/main.go <<EOF
package main

import (
    "github.com/componego/componego"
    "github.com/componego/componego/impl/runner"
    "github.com/componego/componego/libs/color"

    "github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/application"
)

func main() {
    color.SetIsActive(true)
    // This is an entry point for launching the application in developer mode.
    runner.RunAndExit(application.New(), componego.DeveloperMode)
}
EOF

mkdir -p "internal/application"
cat >internal/application/application.go <<EOF
package application

import (
    "fmt"

    "github.com/componego/componego"
    "github.com/componego/componego/impl/application"
)

type Application struct {
}

func New() *Application {
    return &Application{}
}

// ApplicationName belongs to interface componego.Application.
func (a *Application) ApplicationName() string {
    return "Application | v0.0.1"
}

// TODO: you can remove unused methods.

// ApplicationComponents belongs to interface componego.ApplicationComponents.
func (a *Application) ApplicationComponents() ([]componego.Component, error) {
    // Documentation -> https://componego.github.io/impl/application/#applicationcomponents
    return nil, nil
}

// ApplicationDependencies belongs to interface componego.ApplicationDependencies.
func (a *Application) ApplicationDependencies() ([]componego.Dependency, error) {
    // Documentation -> https://componego.github.io/impl/application/#applicationdependencies
    return nil, nil
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *Application) ApplicationConfigInit(appMode componego.ApplicationMode) (map[string]any, error) {
    // Documentation -> https://componego.github.io/impl/application/#applicationconfiginit
    return nil, nil
}

// ApplicationErrorHandler belongs to interface componego.ApplicationErrorHandler.
func (a *Application) ApplicationErrorHandler(err error, appIO componego.ApplicationIO, appMode componego.ApplicationMode) error {
    // Documentation -> https://componego.github.io/impl/application/#applicationerrorhandler
    return err
}

// ApplicationAction belongs to interface componego.Application.
func (a *Application) ApplicationAction(env componego.Environment, args []string) (int, error) {
    // Documentation -> https://componego.github.io/impl/application/#applicationaction
    // In this function you can write business logic for your application.

    // TODO: implement me.

    _, err := fmt.Fprintln(env.ApplicationIO().OutputWriter(), "Hello.")
    return application.ExitWrapper(err)
}

var (
    _ componego.Application             = (*Application)(nil)
    _ componego.ApplicationComponents   = (*Application)(nil)
    _ componego.ApplicationDependencies = (*Application)(nil)
    _ componego.ApplicationConfigInit   = (*Application)(nil)
    _ componego.ApplicationErrorHandler = (*Application)(nil)
)
EOF

mkdir -p "tests"
cat >tests/basic_test.go <<EOF
package tests

import (
    "testing"

    "github.com/componego/componego/tests/runner"

    "github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tests/mocks"
)

func TestBasic(t *testing.T) {
    // Documentation -> https://componego.github.io/tests/runner/
    env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock())
    t.Cleanup(cancelEnv)

    // TODO: implement tests.
    _ = env
}
EOF

mkdir -p "tests/mocks"
cat >tests/mocks/application.go <<EOF
package mocks

import (
    "github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/application"
)

type ApplicationMock struct {
    *application.Application
}

func NewApplicationMock() *ApplicationMock {
    return &ApplicationMock{
        Application: application.New(),
    }
}

// Documentation -> https://componego.github.io/tests/mock/
// ... your other methods
EOF

cat >README.md <<EOF
# Hello

Thank you for creating a basic version of the project based on our framework.

Please visit the [documentation pages](https://componego.github.io/) to understand how to use it.

---

To build the application in developer mode, use the following file: [cmd/application/dev/main.go](./cmd/application/dev/main.go)

There is another file for production mode: [cmd/application/main.go](./cmd/application/main.go)
EOF
