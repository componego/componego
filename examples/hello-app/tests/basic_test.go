package tests

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/tests/runner"

	"github.com/componego/componego/examples/hello-app/tests/mocks"
)

func TestBasic(t *testing.T) {
	buffer := &bytes.Buffer{}
	// We run tests inside mock of the current application example.
	// You can replace parts of the application specifically for the test in this application mock.
	env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock(), &runner.TestOptions{
		Driver: driver.New(&driver.Options{
			AppIO: application.NewIO(nil, buffer, buffer),
		}),
	})
	t.Cleanup(cancelEnv)

	exitCode, err := env.Application().ApplicationAction(env, nil)
	if exitCode != componego.SuccessExitCode || err != nil {
		t.Fatal("the application stopped with an error")
	}
	if buffer.String() != fmt.Sprintln("Hello World!") {
		t.Fatal("different application output was expected")
	}
}
