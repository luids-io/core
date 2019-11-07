// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

// Service the interface implemented by microservice stubs
type Service interface {
	API() string
	Close() error
	Ping() error
}

// Discover interface defines methods for service discovering
type Discover interface {
	GetService(id string) (Service, bool)
	ListServices() []string
}
