// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist

import "errors"

// Some standard errors returned by check interfaces
var (
	ErrResourceNotSupported = errors.New("resource is not supported")
	ErrBadResourceFormat    = errors.New("bad format in request")
	ErrListNotAvailable     = errors.New("list is not available")
)
