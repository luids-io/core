// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist

import "errors"

// Some standard errors returned by List interfaces
var (
	ErrCanceledRequest = errors.New("canceled request")
	ErrBadRequest      = errors.New("bad request")
	ErrNotSupported    = errors.New("resource not supported")
	ErrUnavailable     = errors.New("not available")
	ErrInternal        = errors.New("internal error")
)
