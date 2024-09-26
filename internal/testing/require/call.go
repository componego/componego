/*
Copyright 2024-present Volodymyr Konstanchuk and contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package require

import (
	"errors"

	"github.com/componego/componego/internal/testing"
	"github.com/componego/componego/libs/vendor-proxy"
)

func call(name string, t testing.T, args ...any) {
	if h, ok := t.(testing.THelper); ok {
		h.Helper() // Removes this function from the error stack.
	}
	var err error
	require := vendor_proxy.Get("testify/require")
	if len(args) == 0 {
		_, err = require.CallFunction(name, t)
	} else {
		newArgs := make([]any, 0, len(args)+1)
		newArgs = append(newArgs, t)
		newArgs = append(newArgs, args...)
		_, err = require.CallFunction(name, newArgs...)
	}
	if err == nil {
		return
	}
	message := "error: " + err.Error()
	if errors.Is(err, vendor_proxy.ErrFunctionNotExists) {
		message += messageOfIncorrectRun()
	}
	t.Errorf(message)
	t.FailNow()
}

func messageOfIncorrectRun() string {
	// noinspection SpellCheckingInspection
	return `
Make sure you run tests using 'make tests' or 'make tests-cover'.
However, if you are running tests in your application, then you need to add a new imports:
<code>
    import (
        _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-app"
        _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-tests"
    )
</code>
`
}
