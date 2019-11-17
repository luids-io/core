// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package stream provides a hyperscan stream client and a ready-to-use service
// wrapper for the luids.hyperscan.v1.Block API
//
// This package is a work in progress and makes no API stability promises.
package stream

import "fmt"

// Constants for api description
const (
	APIName    = "luids.hyperscan"
	APIVersion = "v1"
	APIService = "Stream"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
