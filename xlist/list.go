// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package xlist

import (
	"context"
)

// List is the main interface for RBL lists
type List interface {
	Writer
	Checker
}

// Writer is the interface for write in lists
type Writer interface {
	// Append to the list a resource (or group) encoded as string in the format
	Append(ctx context.Context, name string, r Resource, f Format) error
	// Remove from the list
	Remove(ctx context.Context, name string, r Resource, f Format) error
	// Clear all items in the list
	Clear(ctx context.Context) error
	// ReadOnly returns true if the list is read only
	ReadOnly() (bool, error)
}

// Checker is the interface for check lists
type Checker interface {
	// Check method checks if the value encoded as string is in the list
	Check(ctx context.Context, name string, r Resource) (Response, error)
	// Resources returns an array with the resource types supported by the RBL service
	Resources() []Resource
	// Ping method allows to check if the RBL service is working
	Ping() error
}

//Response stores information about the service's responses
type Response struct {
	// Result is true if the value is in the list
	Result bool `json:"result"`
	// Reason stores the reason why it is the list (or not if you want)
	Reason string `json:"reason,omitempty"`
	// TTL is a number in seconds used for caching
	TTL int `json:"ttl"`
}

// NeverCache is a special value for TTL. If TTLs has this value, caches
// should not store the response
const NeverCache = -1
