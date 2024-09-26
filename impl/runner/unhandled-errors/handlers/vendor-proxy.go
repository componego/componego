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

package handlers

import (
	"errors"
	"io"

	"github.com/componego/componego"
	"github.com/componego/componego/internal/utils"
	"github.com/componego/componego/libs/vendor-proxy"
)

func VendorProxyHandler(err error, writer io.Writer, appMode componego.ApplicationMode) bool {
	if !errors.Is(err, vendor_proxy.ErrFunctionNotExists) {
		return false
	}
	DefaultHandler(err, writer, appMode)
	utils.Fprint(writer, `---
Perhaps, you need to add the following lines to your main package.
<code>
    import (
        _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-app"   // if you run the application.
        _ "github.com/componego/meta-package/pre-init/vendor-proxy/for-tests" // if you run the tests.
    )
</code>
If that doesn't help, then you need to debug the code.
`)
	return true
}
