// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package notify

import "fmt"

// Constants for api description
const (
	APIName    = "luids.notify"
	APIVersion = "v1"
	APIService = "Notify"
)

// ServiceName returns service name
func ServiceName() string {
	return fmt.Sprintf("%s.%s.%s", APIName, APIVersion, APIService)
}
