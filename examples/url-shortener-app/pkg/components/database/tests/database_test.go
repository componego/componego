package tests

import (
	"testing"

	"github.com/componego/componego/tests/runner"

	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database"
	"github.com/componego/componego/examples/url-shortener-app/pkg/components/database/tests/mocks"
)

func TestComponent(t *testing.T) {
	env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock())
	t.Cleanup(cancelEnv)
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		_, err := env.DependencyInvoker().Invoke(func(dbProvider database.Provider) error {
			db, err := dbProvider.Get("test-storage")
			if err != nil {
				return err
			}
			var version string
			if err = db.QueryRow(`SELECT VERSION()`).Scan(&version); err != nil {
				return err
			}
			if version != "0.0.1" {
				t.Fatal("no expected data")
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}
