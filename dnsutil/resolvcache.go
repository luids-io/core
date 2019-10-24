// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package dnsutil

import (
	"context"
	"errors"
	"net"
	"time"
)

// Some standard errors returned by resolv interfaces
var (
	ErrBadRequestFormat      = errors.New("bad format in request")
	ErrCacheNotAvailable     = errors.New("cache is not available")
	ErrCollectDNSClientLimit = errors.New("max queries per dns client")
	ErrCollectNamesLimit     = errors.New("max names resolved for an ip")
)

// ResolvCache interface defines a cache for dns resolutions
type ResolvCache interface {
	ResolvCollector
	ResolvChecker
}

// ResolvCollector interface collects information of resolved ips
type ResolvCollector interface {
	Collect(ctx context.Context, client net.IP, name string, resolved []net.IP) error
}

// ResolvChecker is the interface for checks in a resolv cache
type ResolvChecker interface {
	Check(ctx context.Context, client, resolved net.IP, name string) (ResolvResponse, error)
}

// ResolvResponse stores resolv
type ResolvResponse struct {
	// Result is true if was resolved
	Result bool `json:"result"`
	// Last time resolved
	Last time.Time `json:"last,omitempty"`
	// Store time
	Store time.Time `json:"store"`
}
