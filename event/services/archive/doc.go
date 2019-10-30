// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import "fmt"

// Constants for api description
const (
	APIName    = "luids.event"
	APIVersion = "v1"
	APIService = "Archive"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
