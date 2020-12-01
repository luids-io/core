// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/core/grpctls"
)

// ServiceDef is used for define and construct microservices
type ServiceDef struct {
	// ID must exist and be unique for its correct operation
	ID string `json:"id"`
	// Disabled
	Disabled bool `json:"disabled,omitempty"`
	// API defines the api implemented by the service
	API string `json:"api"`
	// Endpoint url
	Endpoint string `json:"endpoint"`
	// Client configuration
	Client *grpctls.ClientCfg `json:"client,omitempty"`
	// Enable log
	Log bool `json:"log,omitempty"`
	// Enable metrics
	Metrics bool `json:"metrics,omitempty"`
	// Enable cache
	Cache bool `json:"cache,omitempty"`
	// Opts stores custom fields
	Opts map[string]interface{} `json:"opts,omitempty"`
}

// Validate checks field values.
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

// ClientCfg returns a copy of client configuration.
// It returns an empty struct if a null pointer is stored.
func (def ServiceDef) ClientCfg() grpctls.ClientCfg {
	if def.Client == nil {
		return grpctls.ClientCfg{}
	}
	return *def.Client
}

// ServiceDefsFromFile reads from file a slice of ServiceDef.
func ServiceDefsFromFile(path string) ([]ServiceDef, error) {
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
