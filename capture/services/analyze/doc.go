// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package analyze

import "fmt"

// Constants for api description
const (
	APIName    = "luids.capture"
	APIVersion = "v1"
	APIService = "Analyze"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
