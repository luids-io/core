// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package resolvcheck implements a resolv checker and a ready
// to use service component.
//
// This package is a work in progress and makes no API stability promises.
package resolvcheck

import "fmt"

// Constants for api description
const (
	APIName    = "luids.dnsutil"
	APIVersion = "v1"
	APIService = "ResolvCheck"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
