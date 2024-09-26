#!/usr/bin/env python

# Copyright 2024-present Volodymyr Konstanchuk and contributors
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

""" A set of utilities for checking code """

# pylint make.py --disable=C0116,C0115,W0511 --max-line-length=120

from sys import version_info, exit, stdout, argv, executable  # pylint: disable=redefined-builtin

if version_info < (3, 10):
    print('This script supports Python 3.10 or later.')
    exit(2)

# pylint: disable=wrong-import-position

from os import path, environ, makedirs, symlink, sep as path_separator, fdopen, close, pipe
from tempfile import TemporaryDirectory
from subprocess import run as runprocess, SubprocessError, PIPE as SUBPROCESS_PIPE, Popen
from glob import iglob, glob
from shutil import copy, move as rename
from typing import TypeAlias, Final, Callable, IO, Any
from threading import Thread
from atexit import register as on_exit
from uuid import uuid4
from hashlib import sha1
from re import search as re_search
from concurrent import futures as concurrent

# pylint: enable=wrong-import-position

META_PACKAGE_VERSION: Final[str] = 'latest'
# noinspection SpellCheckingInspection
GOSEC_VERSION: Final[str] = 'latest'
# noinspection SpellCheckingInspection
GOLANGCI_LINT_VERSION: Final[str] = 'v1.61.0'
# noinspection SpellCheckingInspection
GOFUMPT_VERSION: Final[str] = 'v0.7.0'

# This is the Apache 2.0 license hash that is present in each file.
LICENSE_HASH: Final[str] = 'c5d3aea6439d6869df0a3f8d705df6ff392cd858'

TEST_COVERAGE_LIMIT: Final[float] = 80
# noinspection SpellCheckingInspection
TEST_COVERAGE_IGNORE: Final[tuple] = (
    # all examples.
    'examples/**',
    # all test packages except the root package.
    'impl/**/tests/**',
    'internal/**/tests/**',
    'libs/**/tests/**',
    'tools/**/tests/**',
)
TEST_COVERAGE_FILE: Final[str] = 'coverage.out'

__all__ = []

Args: TypeAlias = tuple[str] | str
File: TypeAlias = int | IO[Any]


class TestEnvironment(TemporaryDirectory):
    # noinspection SpellCheckingInspection
    TESTS_INIT_CODE: Final[str] = """
// dependencies init file

import (
    "fmt"
    "io"

    "github.com/componego/componego/internal/system"

    _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-app"
    _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-tests"
)

func init() {
    // This is callback output, which is responsible for identifying the go test command.
    fmt.Println("[ID]")
    // Disable console output of functions inside the framework.
    system.Stdout = io.Discard
    system.Stderr = io.Discard
}
"""
    _env_id: str

    def __init__(self) -> None:
        # noinspection SpellCheckingInspection
        super().__init__(prefix='componego_')
        self._env_id = 'test:' + uuid4().hex
        self._init_vendor_proxy()

    def __enter__(self):
        return self.name, self._env_id

    def _init_vendor_proxy(self) -> None:
        src_dir = basedir()
        dst_dir = self.name
        copy(path.join(src_dir, 'go.mod'), path.join(dst_dir, 'go.mod'))
        if path.exists(path.join(src_dir, 'go.sum')):
            copy(path.join(src_dir, 'go.sum'), path.join(dst_dir, 'go.sum'))
        # noinspection SpellCheckingInspection
        require_path = path.join('internal', 'testing', 'require')
        makedirs(path.join(dst_dir, require_path))
        call_file = path.join(src_dir, require_path, 'call.go')
        call_init_file = path.join(dst_dir, require_path, '0.go')
        self._write_init_file(call_file, call_init_file)
        self._create_symlinks(src_dir, dst_dir)
        # noinspection SpellCheckingInspection
        Go.get('github.com/componego/meta-package', version=META_PACKAGE_VERSION, cwd=dst_dir)
        Go.tidy(dst_dir)

    def _write_init_file(self, src_file: str, dst_file: str) -> None:
        with open(src_file, 'r', encoding='utf-8') as reader, open(dst_file, 'w', encoding='utf-8') as writer:
            for line in reader:
                writer.write(line)
                if not line.startswith('package '):
                    continue
                writer.write(self.TESTS_INIT_CODE.replace('[ID]', self._env_id))
                break

    @staticmethod
    def _create_symlinks(src_dir: str, dst_dir: str) -> None:
        for src_file in iglob('**/*_test.go', root_dir=src_dir, recursive=True):
            dst_file = path.join(dst_dir, src_file)
            src_file = path.join(src_dir, src_file)
            makedirs(path.dirname(dst_file), exist_ok=True)
            symlink(src_file, dst_file)
        for src_file in iglob('**/*.go', root_dir=src_dir, recursive=True):
            src_file_items = filter(None, src_file.split(path_separator))
            dst_file = dst_dir
            for item in src_file_items:
                dst_file = path.join(dst_file, item)
                if path.exists(dst_file):
                    continue
                symlink(dst_file.replace(dst_dir, src_dir, 1), dst_file)


class Go:
    @staticmethod
    def root() -> str:
        # noinspection SpellCheckingInspection
        key = 'GOROOT'
        root = environ.get(key, None)
        if root is None:
            info = runprocess(f'go env {key}', capture_output=True, text=True, shell=True, check=True)
            root = str(info.stdout).strip()
            environ[key] = root
        if len(root) == 0:
            raise ValueError(f'missing {key} environment variable')
        return root.split(':', 1)[0]

    @classmethod
    def bin(cls, name: str = None) -> str:
        result = path.join(cls.root(), 'bin')
        if name is None:
            return result
        return path.join(result, name)

    @classmethod
    def install(cls, package: str, version: str) -> None:
        binary_name = path.basename(package)
        if path.exists(cls.bin(binary_name)):
            return
        cached_file = path.join(path.expanduser('~'), '.cache', 'componego', 'bin', binary_name)
        if is_ci_cd_pipeline() and path.exists(cached_file):
            print(f'install from cache: {package}')
            copy(cached_file, cls.bin(binary_name))
            return
        package = f'{package}@{version}'
        print(f'install: {package}')
        # noinspection SpellCheckingInspection
        with TemporaryDirectory(prefix='componego_') as tempdir:
            # noinspection SpellCheckingInspection
            env = environ | {'GOBIN': cls.bin()}
            run_process('go install', args=package, cwd=tempdir, env=env)
        if is_ci_cd_pipeline():
            makedirs(path.dirname(cached_file), exist_ok=True)
            copy(cls.bin(binary_name), cached_file)

    @classmethod
    def get(cls, package: str, version: str, cwd: str = None) -> None:
        run_process(f'go get {package}@{version}', cwd=cwd)

    @classmethod
    def tidy(cls, cwd: str) -> None:
        run_process('go mod tidy', cwd=cwd)

    @classmethod
    def run(cls, args: Args, cwd: str) -> None:
        run_process('go run', args=args, cwd=cwd)


class Command:
    _instances: dict[str, Callable] = {}

    def __init__(self, function: Callable[[Args], None]) -> None:
        self._function = function
        self._register(function.__name__, self)

    def __call__(self, *args, **kwargs) -> None:
        name = self._function.__name__.replace('_', ':')
        print(f'> RUN > {name}')
        self._function(*args, **kwargs)

    @staticmethod
    def args(current: Args, default: Args) -> Args:
        return current if len(current) > 0 else default

    @classmethod
    def _register(cls, name: str, instance: Callable) -> None:
        cls._instances[name.replace('_', ':')] = instance

    @classmethod
    def main(cls) -> int:
        try:
            if len(argv) >= 2 and argv[1] in cls._instances:
                cls._instances[argv[1]](tuple(argv[2:]))
                return 0
            commands = ', '.join(cls._instances.keys())
            raise ValueError(f'command is missing. The list of allowed commands is {commands}')
        except KeyboardInterrupt:
            print('keyboard interrupt')
        except (Exception,) as e:  # pylint: disable=W0718
            print(f'Error > {e}')
        return 1


def basedir() -> str:
    return path.dirname(path.dirname(path.realpath(__file__)))


def run_tests(cmd: str, args: Args | None, src_dir: str, dst_dir: str, env_id: str) -> None:
    read, write = pipe()

    def pipe_reader():
        reader = fdopen(read)
        can_replace = False
        for line in iter(reader.readline, ''):
            if env_id in line:
                can_replace = True
                continue
            print(line.replace(dst_dir, src_dir) if can_replace else line, end='')
        reader.close()

    env = environ | {
        'CGO_ENABLED': '1',  # for -race flag
    }
    try:
        Thread(target=pipe_reader).start()
        run_process(cmd, args=args, cwd=dst_dir, env=env, output=write)
    finally:
        close(write)


def run_process(cmd: str, args: Args = None, cwd: str = None, env: dict[str, str] = None, output: File = None) -> None:
    if args is not None:
        cmd += ' ' + (args if isinstance(args, str) else ' '.join(args))
    print(f'run command: {cmd}')
    cwd = basedir() if cwd is None else cwd
    output = stdout if output is None else output
    runprocess(cmd, stdout=output, stderr=output, cwd=cwd, env=env, universal_newlines=True, shell=True, check=True)


def is_ci_cd_pipeline() -> bool:
    keys = ['GITHUB_ACTIONS', 'TRAVIS', 'CIRCLECI', 'GITLAB_CI']
    for key in keys:
        if environ.get(key, None):
            return True
    return False


def get_package_name(cwd: str = None) -> str:
    cwd = basedir() if cwd is None else cwd
    with open(path.join(cwd, 'go.mod'), 'r', encoding='utf-8') as file:
        line = file.readline().strip()
        return line.removeprefix('module ').strip()


def ignore_coverage(profile: str, package_dir: str) -> None:
    package_name = get_package_name(package_dir)
    ignore = []
    for ignore_pattern in TEST_COVERAGE_IGNORE:
        ignore.extend(glob(ignore_pattern, root_dir=package_dir, recursive=True))
    tmp_profile = path.join(path.dirname(profile), f'tmp_{path.basename(profile)}')
    with open(profile, 'r', encoding='utf-8') as reader, open(tmp_profile, 'w', encoding='utf-8') as writer:
        for line in reader:
            try:
                filename = line[0:line.find(':')].removeprefix(package_name + '/')
            except ValueError as e:
                raise ValueError('unexpected value in file') from e
            if filename not in ignore:
                writer.write(line)
    rename(tmp_profile, profile)


def upload_coverage(profile: str) -> None:
    is_master_branch = environ.get('GITHUB_REF_NAME') == 'master' and environ.get('GITHUB_REF_TYPE') == 'branch'
    # noinspection SpellCheckingInspection
    is_original_repo = environ.get('GITHUB_REPOSITORY') == 'componego/componego'
    if is_ci_cd_pipeline() and is_master_branch and is_original_repo:
        # This file will be processed in the next step of GitHub actions.
        copy(profile, path.join(basedir(), TEST_COVERAGE_FILE))


def run_parallel_commands(commands: list[str]) -> None:
    def run_command(command: str) -> tuple[str, str, int]:
        print('> run in parallel >', command)
        command = [executable, path.realpath(__file__), command]
        with Popen(command, stdout=SUBPROCESS_PIPE, stderr=SUBPROCESS_PIPE) as process:
            process_stdout, process_stderr = process.communicate()
            return process_stdout.decode('utf-8'), process_stderr.decode('utf-8'), process.returncode

    results: dict[str, tuple[str, str, int]] = {}
    with concurrent.ThreadPoolExecutor(max_workers=len(commands)) as executor:
        futures = {executor.submit(run_command, command): command for command in commands}
        print('> wait for results...')
        for future in futures:
            results[futures[future]] = future.result()

    for command in commands:
        if command not in results:
            continue
        result = results[command]
        print(result[0])
        if result[1]:
            print(result[1])
        if result[2] != 0:
            raise SubprocessError(f'exit code: {result[2]}')


@Command
def fmt(_: Args) -> None:
    run_process('go fmt ./...')
    # noinspection SpellCheckingInspection
    Go.install('mvdan.cc/gofumpt', version=GOFUMPT_VERSION)
    # noinspection SpellCheckingInspection
    run_process(Go.bin('gofumpt'), args='-l -w .')


@Command
def tests(args: Args) -> None:
    # noinspection SpellCheckingInspection
    args = Command.args(args, '-race -v -count=1 ./... -bench=. -benchtime=1x')
    with TestEnvironment() as (tempdir, env_id):
        run_tests('go test', args=args, src_dir=basedir(), dst_dir=tempdir, env_id=env_id)


@Command
def tests_cover(args: Args) -> None:
    percentage = 0
    # noinspection SpellCheckingInspection
    args = Command.args(args, f'-race -v -count=1 -coverpkg=./... -coverprofile={TEST_COVERAGE_FILE} ./...')
    with TestEnvironment() as (tempdir, env_id):
        run_tests('go test', args=args, src_dir=basedir(), dst_dir=tempdir, env_id=env_id)
        # The file may not exist depending on the function arguments.
        if not path.exists(path.join(tempdir, TEST_COVERAGE_FILE)):
            return
        # Some files are ignored when checking test coverage.
        # For example, we ignore the folder with examples, testers and mocks, which are used to run tests.
        # See the function implementation for details.
        ignore_coverage(path.join(tempdir, TEST_COVERAGE_FILE), tempdir)
        upload_coverage(path.join(tempdir, TEST_COVERAGE_FILE))
        output = path.join(tempdir, 'tool.out')
        command = f'go tool cover -func={TEST_COVERAGE_FILE} > {output} 2>&1'
        run_tests(command, args=None, src_dir=basedir(), dst_dir=tempdir, env_id=env_id)
        with open(output, 'r', encoding='utf-8') as reader:
            for line in reader:
                print(line, end='')
                if not line.startswith('total:'):
                    continue
                match = re_search(r'(\d+\.\d+)%', line)
                if match:
                    percentage = float(match.group(1))
                else:
                    raise ValueError('no percentage found in the input string')
    if percentage < TEST_COVERAGE_LIMIT:
        raise ValueError(f'test coverage less than {TEST_COVERAGE_LIMIT}%')


@Command
def tests_env(_: Args) -> None:
    try:
        import readline  # pylint: disable=import-outside-toplevel
    except ImportError:
        readline = None
    if readline is not None:
        readline.set_history_length(100)
        # noinspection SpellCheckingInspection
        history_file = path.expanduser('~/.componego_dev_tests_history')
        if path.exists(history_file):
            readline.read_history_file(history_file)
        on_exit(readline.write_history_file, history_file)
    while True:
        with TestEnvironment() as (tempdir, env_id):
            print(f'test environment is initialized - {tempdir}')
            while True:
                command = input('>>> ')
                if command == 'exit':
                    return
                if command == 'reinit':
                    break
                if not command.strip():
                    continue
                try:
                    run_tests(command, args=None, src_dir=basedir(), dst_dir=tempdir, env_id=env_id)
                except KeyboardInterrupt:
                    print('keyboard interrupt')
                except SubprocessError as e:
                    print(f'Error > {e}')


@Command
def test_skeleton(_: Args) -> None:
    # noinspection SpellCheckingInspection
    with TemporaryDirectory(prefix='componego_') as tempdir:
        copy(path.join(basedir(), 'tools', 'create-basic-app.sh'), path.join(tempdir, 'create-basic-app.sh'))
        run_process('sh create-basic-app.sh', cwd=tempdir)


@Command
def lint(args: Args) -> None:
    args = Command.args(args, 'run ./...')
    # noinspection SpellCheckingInspection
    Go.install('github.com/golangci/golangci-lint/cmd/golangci-lint', version=GOLANGCI_LINT_VERSION)
    # noinspection SpellCheckingInspection
    run_process(Go.bin('golangci-lint'), args=args)


@Command
def security(args: Args) -> None:
    args = Command.args(args, './...')
    # noinspection SpellCheckingInspection
    Go.install('github.com/securego/gosec/v2/cmd/gosec', version=GOSEC_VERSION)
    # noinspection SpellCheckingInspection
    run_process(Go.bin('gosec'), args=args)


@Command
def generate(_: Args) -> None:
    for filename in iglob('**/*.go', root_dir=basedir(), recursive=True):
        filename = path.join(basedir(), filename)
        with open(filename, 'r', encoding='utf-8') as file:
            if 'go:generate' in file.read():
                run_process(f'go generate {filename}', cwd=basedir())


@Command
def validate_dependencies(_: Args) -> None:
    Go.tidy(basedir())
    # noinspection SpellCheckingInspection
    gosum = path.join(basedir(), 'go.sum')
    if not path.exists(gosum):
        raise FileNotFoundError('go.sum not found in the repository')
    if path.getsize(gosum) > 0:
        # noinspection SpellCheckingInspection
        raise ValueError('componego has dependencies but shouldn\'t')


@Command
def license_check(_: Args) -> None:
    """
    This function checks that there is a license at the top of each Go file.
    """
    for filename in iglob('**/*.go', root_dir=basedir(), recursive=True):
        if filename.startswith(f'examples{path_separator}') or filename.startswith(f'docs{path_separator}'):
            continue
        filename = path.join(basedir(), filename)
        with open(filename, 'r', encoding='utf-8') as file:
            text = file.read(700)
            has_copyright = False
            license_text = []
            try:
                for line in text[text.find('/*') + len('/*'):text.rfind('*/')].splitlines():
                    if not has_copyright and line.startswith('Copyright'):
                        has_copyright = True
                    elif has_copyright:
                        license_text.append(line)
                license_hash = sha1('\n'.join(license_text).strip().encode('utf-8')).hexdigest()
                if license_hash == LICENSE_HASH:
                    continue
            except ValueError:
                pass
        # Please check the other files and insert the correct license text at the beginning of the file.
        raise ValueError(f'you should not change the license text at the beginning of each file -> {filename}')


@Command
def github_actions(_: Args) -> None:
    no_args = ''
    validate_dependencies(no_args)
    generate(no_args)
    fmt(no_args)
    try:
        # return an error if there is a difference in the code after generation
        run_process('git diff --quiet', cwd=basedir())
    except SubprocessError as e:
        run_process('git --no-pager diff', cwd=basedir())
        raise e
    run_parallel_commands(['tests', 'test:skeleton', 'security', 'lint', 'license:check', 'tests:cover'])


@Command
def commit_hook(_: Args) -> None:
    no_args = ''
    github_actions(no_args)


if __name__ == '__main__':
    exit(Command.main())
