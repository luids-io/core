// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package apiservice provides a simple system to instantiate clients of the
// luIDS api.
//
// This package is a work in progress and makes no API stability promises.
package apiservice

// Service is the interface that must be implemented by service API clients.
type Service interface {
	// API returns signature
	API() string
	// Close client
	Close() error
	// Ping status
	Ping() error
}

// Discover interface used for service discovering.
type Discover interface {
	// GetService by id, returning false if not available
	GetService(id string) (Service, bool)
	// ListServices returns the ids of services available
	ListServices() []string
}
