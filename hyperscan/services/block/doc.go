// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package block provides a hyperscan block client and a ready-to-use service
// wrapper for the luids.hyperscan.v1.Block API
//
// This package is a work in progress and makes no API stability promises.
package block

import "fmt"

// Constants for api description
const (
	APIName    = "luids.hyperscan"
	APIVersion = "v1"
	APIService = "Block"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
