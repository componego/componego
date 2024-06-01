/*
Copyright 2024 Volodymyr Konstanchuk and the Componego Framework contributors

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

package debug

import (
	"fmt"
	"runtime"
	"strings"
)

type StackTrace []uintptr

func GetStackTrace() *StackTrace {
	const depth = 32
	var pcs [depth]uintptr
	position := runtime.Callers(3, pcs[:])
	var stackTrace StackTrace = pcs[0:position]
	return &stackTrace
}

func (s *StackTrace) String() string {
	var builder strings.Builder
	for _, pc := range *s {
		frame := StackFrame(pc)
		builder.WriteString(frame.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

type StackFrame uintptr

func (s StackFrame) PC() uintptr {
	return uintptr(s) - 1
}

func (s StackFrame) FileLine() (string, int) {
	fn := runtime.FuncForPC(s.PC())
	if fn == nil {
		return "unknown", 0
	}
	return fn.FileLine(s.PC())
}

func (s StackFrame) Name() string {
	fn := runtime.FuncForPC(s.PC())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func (s StackFrame) String() string {
	file, line := s.FileLine()
	return fmt.Sprintf("%s: %s:%d", s.Name(), file, line)
}

var (
	_ fmt.Stringer = (*StackTrace)(nil)
	_ fmt.Stringer = (*StackFrame)(nil)
)
