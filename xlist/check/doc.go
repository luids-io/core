// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package check provides an xlist client and a ready-to-use service
// wrapper for the luids.xlist.vX.Check API
//
// This package is a work in progress and makes no API stability promises.
package check

import "fmt"

// Constants for api description
const (
	APIName    = "luids.xlist"
	APIVersion = "v1"
	APIService = "Check"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
