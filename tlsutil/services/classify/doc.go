// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package classify provides a machine learning classifier client and a
// ready-to-use service wrapper for the API
//
// This package is a work in progress and makes no API stability promises.
package classify

import "fmt"

// Constants for api description
const (
	APIName    = "luids.tlsutil"
	APIVersion = "v1"
	APIService = "Classify"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
