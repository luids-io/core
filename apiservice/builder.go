// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisguillenc/grpctls"
	"github.com/luisguillenc/yalogi"
)

// ServiceDef is used for define and construct microservices
type ServiceDef struct {
	// ID must exist and be unique for its correct operation
	ID string `json:"id"`
	// API defines the api implemented by the service
	API string `json:"api"`
	// Endpoint url
	Endpoint string `json:"endpoint"`
	// Client configuration
	Client *grpctls.ClientCfg `json:"client,omitempty"`
	// Metrics
	Metrics bool `json:"metrics,omitempty"`
	// Disabled
	Disabled bool `json:"disabled,omitempty"`
	// Opts allow custom fields
	Opts map[string]interface{} `json:"opts,omitempty"`
}

// BuildFn defines a function that constructs a service using a
// definition
type BuildFn func(def ServiceDef, logger yalogi.Logger) (Service, error)

// Validate checks definition field values
func (def ServiceDef) Validate() error {
	if def.ID == "" {
		return errors.New("'id' is required")
	}
	if def.API == "" {
		return errors.New("'api' is required")
	}
	//parses endpoint
	_, _, err := grpctls.ParseURI(def.Endpoint)
	if err != nil {
		return fmt.Errorf("'endpoint' invalid: %v", err)
	}
	//grpc client config
	if def.Client != nil {
		err = def.Client.Validate()
		if err != nil {
			return fmt.Errorf("'config' invalid: %v", err)
		}
	}
	return nil
}

// ClientCfg returns client configuration
func (def ServiceDef) ClientCfg() grpctls.ClientCfg {
	if def.Client == nil {
		return grpctls.ClientCfg{}
	}
	return *def.Client
}

// RegisterBuilder registers a service builder for an api signature
func RegisterBuilder(api string, builder BuildFn) {
	registryBuilder[api] = builder
}

// Build creates a new service using a service definition struct
func Build(def ServiceDef, logger yalogi.Logger) (Service, error) {
	if def.Disabled {
		return nil, errors.New("service is disabled")
	}
	if def.API == "" {
		return nil, errors.New("'api' is required")
	}
	//get builder for related api
	customb, ok := registryBuilder[def.API]
	if !ok {
		return nil, errors.New("'api' not registered")
	}
	return customb(def, logger)
}

// DefsFromFile creates a slice of Definition from a file in json format.
func DefsFromFile(path string) ([]ServiceDef, error) {
	var services []ServiceDef
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byteValue, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}

var registryBuilder map[string]BuildFn

func init() {
	registryBuilder = make(map[string]BuildFn)
}
