// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist

import "errors"

// Some standard errors returned by List interfaces
var (
	ErrNotImplemented = errors.New("not implemented")
	ErrBadRequest     = errors.New("bad request")
	ErrNotAvailable   = errors.New("not available")
	ErrReadOnlyMode   = errors.New("read only mode")
)
