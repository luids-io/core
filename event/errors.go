// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package event

import "errors"

// Some standard errors returned by interfaces
var (
	ErrCanceledRequest = errors.New("canceled request")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotSupported    = errors.New("not supported")
	ErrUnavailable     = errors.New("not available")
	ErrInternal        = errors.New("internal error")
)
