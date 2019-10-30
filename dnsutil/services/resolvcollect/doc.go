// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package resolvcollect implements a resolv collector client and a ready
// to use service component.
//
// This package is a work in progress and makes no API stability promises.
package resolvcollect

import "fmt"

// Constants for api description
const (
	APIName    = "luids.dnsutil"
	APIVersion = "v1"
	APIService = "ResolvCollect"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
