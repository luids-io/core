// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package resolvarchive implements a resolv archiver client and a ready
// to use service component.
//
// This package is a work in progress and makes no API stability promises.
package resolvarchive

import "fmt"

// Constants for api description
const (
	APIName    = "luids.dnsutil"
	APIVersion = "v1"
	APIService = "ResolvArchive"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
