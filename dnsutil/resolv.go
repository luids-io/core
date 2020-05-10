// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package dnsutil

import (
	"context"
	"net"
	"time"
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
