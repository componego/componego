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

package xerrors

import (
	"errors"
)

type XError interface {
	error
	ErrorCode() string
	ErrorOptions() []Option
	WithMessage(message string, code string, options ...Option) XError
	WithError(err error, code string, options ...Option) XError
	WithOptions(code string, options ...Option) XError
}

type xError struct {
	parent    XError
	current   error
	asMessage bool
	code      string
	options   []Option
}

func New(message string, code string, options ...Option) XError {
	return &xError{
		parent:    nil,
		current:   errors.New(message),
		asMessage: true,
		code:      code,
		options:   options,
	}
}

func ToXError(err error, code string, options ...Option) XError {
	return &xError{
		parent:    nil,
		current:   err,
		asMessage: false,
		code:      code,
		options:   options,
	}
}

func (x *xError) Error() string {
	if x.parent == nil {
		if x.code == "" {
			return x.current.Error()
		}
		return x.current.Error() + " (" + x.code + ")"
	}
	if x.code == "" {
		return x.parent.Error() + "\n" + x.current.Error()
	}
	return x.parent.Error() + "\n" + x.current.Error() + " (" + x.code + ")"
}

func (x *xError) ErrorCode() string {
	return x.code
}

func (x *xError) ErrorOptions() []Option {
	return x.options
}

func (x *xError) WithMessage(message string, code string, options ...Option) XError {
	return &xError{
		parent:    x,
		current:   errors.New(message),
		asMessage: true,
		code:      code,
		options:   options,
	}
}

func (x *xError) WithError(err error, code string, options ...Option) XError {
	return &xError{
		parent:    x,
		current:   err,
		asMessage: false,
		code:      code,
		options:   options,
	}
}

func (x *xError) WithOptions(code string, options ...Option) XError {
	return &xError{
		parent:    nil,
		current:   x,
		asMessage: false,
		code:      code,
		options:   options,
	}
}

func (x *xError) Unwrap() []error {
	if x.parent == nil {
		if x.asMessage {
			return nil
		}
		return []error{x.current}
	}
	if x.asMessage {
		return []error{x.parent}
	}
	return []error{x.parent, x.current}
}

var (
	_ error                         = (*xError)(nil)
	_ XError                        = (*xError)(nil)
	_ interface{ Unwrap() []error } = (*xError)(nil)
)
