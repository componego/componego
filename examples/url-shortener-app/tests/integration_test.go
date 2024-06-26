package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/componego/componego/tests/runner"

	"github.com/componego/componego/examples/url-shortener-app/internal/server"
	"github.com/componego/componego/examples/url-shortener-app/internal/utils"
	"github.com/componego/componego/examples/url-shortener-app/tests/mocks"
)

func TestIntegration(t *testing.T) {
	// We run tests inside mock of the current application example.
	// You can replace parts of the application specifically for the test in this application mock.
	env, cancelEnv := runner.CreateTestEnvironment(t, mocks.NewApplicationMock(), nil)
	t.Cleanup(cancelEnv)
	t.Run("create urls", func(t *testing.T) {
		t.Parallel() // Parallel running of tests is supported.
		router, err := server.CreateRouter(env)
		if err != nil {
			t.Fatalf("create router error: %s", err)
		}
		testServer := httptest.NewServer(router)
		defer testServer.Close()
		longUrl := fmt.Sprintf("https://%s.com/", utils.GetRandomString(100))
		shortUrl := getShortUrl(t, testServer.URL+"/create", longUrl)
		if getLongUrl(t, shortUrl) != longUrl {
			t.Fatal("short and long urls do not match")
		}
	})
}

func getShortUrl(t *testing.T, endpoint string, longUrl string) string {
	response, err := (&http.Client{}).Do((func(body string) *http.Request {
		request, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer([]byte(body)))
		if err != nil {
			t.Fatalf("send request error: %s", err)
		}
		request.Header.Set("Content-Type", "application/json")
		return request
	})(fmt.Sprintf(`{ "url": "%s" }`, longUrl)))
	if err != nil {
		t.Fatalf("send request error: %s", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid response status: %d", response.StatusCode)
	}
	var responseAsStruct struct {
		Status bool   `json:"status"`
		Error  string `json:"error,omitempty"`
		Data   struct {
			NewUrl string `json:"newUrl"`
		} `json:"data,omitempty"`
	}
	if err = json.NewDecoder(response.Body).Decode(&responseAsStruct); err != nil {
		t.Fatal("invalid response received")
	} else if responseAsStruct.Status != true {
		t.Fatal("redirect was not created")
	}
	return responseAsStruct.Data.NewUrl
}

func getLongUrl(t *testing.T, endpoint string) string {
	response, err := (&http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}).Get(endpoint)
	if err != nil {
		t.Fatalf("send request error: %s", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusPermanentRedirect {
		t.Fatal("failed to get redirect")
	}
	return response.Header.Get("Location")
}
