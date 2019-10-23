// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Registry stores service items
type Registry struct {
	services map[string]serviceItem
	mu       sync.RWMutex
}

type serviceItem struct {
	id      string
	service Service
}

// NewRegistry instantiates a new registry
func NewRegistry() *Registry {
	return &Registry{services: make(map[string]serviceItem)}
}

// Register a service using an id and a Service interface
func (r *Registry) Register(id string, svc Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.services[id]
	if ok {
		return errors.New("service already exists")
	}
	r.services[id] = serviceItem{id: id, service: svc}
	return nil
}

// GetService using its id, returns service interface, api signature and true if service exists
func (r *Registry) GetService(id string) (Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.services[id]
	if !ok {
		return nil, false
	}
	return i.service, true
}

// List returns an ordered list of registered ids
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]string, 0, len(r.services))
	for k := range r.services {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// Ping all registered services
func (r *Registry) Ping() error {
	errs := make([]string, 0, len(r.services))
	for _, id := range r.List() {
		svc, ok := r.services[id]
		if ok {
			err := svc.service.Ping()
			errs = append(errs, fmt.Sprintf("%s: %s", svc.id, err.Error()))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}

// CloseAll registered services
func (r *Registry) CloseAll() error {
	errs := make([]string, 0, len(r.services))
	for _, id := range r.List() {
		svc, ok := r.services[id]
		if ok {
			err := svc.service.Close()
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}
