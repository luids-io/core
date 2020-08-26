// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"errors"

	"github.com/luids-io/core/yalogi"
)

// BuildFn defines a function that constructs a service using a service definition.
type BuildFn func(def ServiceDef, logger yalogi.Logger) (Service, error)

// RegisterBuilder registers a service builder for an api signature.
func RegisterBuilder(api string, builder BuildFn) {
	registryBuilder[api] = builder
}

// Build creates a new service using a service definition struct.
func Build(def ServiceDef, logger yalogi.Logger) (Service, error) {
	if def.Disabled {
		return nil, errors.New("apiservice: service is disabled")
	}
	if def.API == "" {
		return nil, errors.New("apiservice: 'api' is required")
	}
	//get builder for related api
	customb, ok := registryBuilder[def.API]
	if !ok {
		return nil, errors.New("apiservice: 'api' not registered")
	}
	return customb(def, logger)
}

// stores builders indexed by api signature
var registryBuilder map[string]BuildFn

func init() {
	registryBuilder = make(map[string]BuildFn)
}
