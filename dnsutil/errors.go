// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package dnsutil

import "errors"

// Some standard errors returned by interfaces
var (
	ErrCanceledRequest = errors.New("canceled request")
	ErrBadRequest      = errors.New("bad request")
	ErrNotSupported    = errors.New("not supported")
	ErrUnavailable     = errors.New("not available")
	ErrInternal        = errors.New("internal error")
	//limit errors
	ErrLimitDNSClientQueries = errors.New("max queries per dns client")
	ErrLimitResolvedNamesIP  = errors.New("max names resolved for an ip")
)
