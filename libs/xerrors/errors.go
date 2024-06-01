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
	Options() []Option
	WithMessage(message string, options ...Option) XError
	WithError(err error, options ...Option) XError
	WithOptions(options ...Option) XError
}

type xError struct {
	parent    XError
	current   error
	asMessage bool
	options   []Option
}

func New(message string, options ...Option) XError {
	return &xError{
		parent:    nil,
		current:   errors.New(message),
		asMessage: true,
		options:   options,
	}
}

func ConvertToXError(err error, options ...Option) XError {
	// noinspection ALL
	if xErr, ok := err.(XError); ok { //nolint:errorlint
		return xErr.WithOptions(options...)
	}
	return &xError{
		parent:    nil,
		current:   err,
		asMessage: false,
		options:   options,
	}
}

func (x *xError) Error() string {
	if x.parent == nil {
		return x.current.Error()
	}
	return x.parent.Error() + " -> " + x.current.Error()
}

func (x *xError) Options() []Option {
	return x.options
}

func (x *xError) WithMessage(message string, options ...Option) XError {
	return &xError{
		parent:    x,
		current:   errors.New(message),
		asMessage: true,
		options:   options,
	}
}

func (x *xError) WithError(err error, options ...Option) XError {
	if errors.Is(err, x) {
		return x.WithOptions(options...)
	}
	return &xError{
		parent:    x,
		current:   err,
		asMessage: false,
		options:   options,
	}
}

func (x *xError) WithOptions(options ...Option) XError {
	if len(options) == 0 {
		return x
	}
	return &xError{
		parent:    nil,
		current:   x,
		asMessage: false,
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
