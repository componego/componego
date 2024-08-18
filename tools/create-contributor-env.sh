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

check_directory() {
    if [ -e "$1" ]; then
        echo "ERROR: could not create a project in the current directory because we are inside another project: $1"
        exit 1
    fi
}

directory=$(pwd)
while [ "$directory" != "/" ]; do
    check_directory "$directory/go.mod"
    check_directory "$directory/.git"
    directory=$(dirname "$directory")
done

mkdir -p "componego-contributor-env"
cd "componego-contributor-env"

cat >docker-compose.yml <<EOF
version: '3.8'

services:
  componego-framework:
    build:
      context: .
      dockerfile: src/componego/scripts/Dockerfile
    container_name: componego-framework-container
    working_dir: /go/src/github.com/componego
    volumes:
      - ./src:/go/src/github.com/componego:cached
  componego-framework-docs:
    build:
      context: ./src/componego/docs/
      dockerfile: Dockerfile
    container_name: componego-framework-docs
    working_dir: /docs
    volumes:
      - ./src/componego/docs:/docs:cached
    ports:
      - "8123:8123"
EOF

cat >.gitattributes <<EOF
* text=auto eol=lf
*.bat text eol=crlf
*.cmd text eol=crlf
*.ahk text eol=crlf
EOF

cat >.gitignore <<EOF
# IDE
.idea
**/.idea
.vs-code
**/.vs-code

# MacOS files
.DS_STORE
**/.DS_Store

# Source code
src/
EOF

cat >.editorconfig <<EOF
root = true

[*]
charset = utf-8
end_of_line = lf
indent_size = 4
indent_style = space
insert_final_newline = true
trim_trailing_whitespace = true

[{*.yml,*.yaml}]
indent_size = 2

[{Makefile,go.mod,go.sum}]
indent_style = tab

[LICENSE]
insert_final_newline = false
EOF

mkdir -p "src"
cat >src/README.md <<EOF
This folder contains repositories.
If it is empty, follow the instructions provided after you create this environment.
EOF

COMMAND_COLOR="\033[0;31m"
RESET_COLOR="\033[0m"

echo "The folder structure has been created -> $(pwd)"
echo "In the next step, run the following commands manually:"
echo ">$"
echo ">$ ${COMMAND_COLOR}cd $(pwd)${RESET_COLOR}"
echo ">$ ${COMMAND_COLOR}git clone https://github.com/componego/componego.git src/componego${RESET_COLOR} # replace repo with your fork if your have one"
echo ">$ ${COMMAND_COLOR}docker-compose up componego-framework -d${RESET_COLOR} # start docker container in background"
echo ">$ ${COMMAND_COLOR}docker inspect --format '{{json .State.Running}}' componego-framework-container${RESET_COLOR} # check if docker container is running"
echo ">$ ${COMMAND_COLOR}docker exec -ti componego-framework-container /bin/bash${RESET_COLOR} # open terminal inside docker container"
echo ">$ ${COMMAND_COLOR}cd componego${RESET_COLOR} # change folder inside docker container"
echo ">$ ${COMMAND_COLOR}make tests${RESET_COLOR} # run tests inside docker container"
echo ">$"
