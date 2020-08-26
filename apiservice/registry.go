// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Registry stores service items indexed by an id. Implements Discover interface.
type Registry struct {
	list     []string
	services map[string]Service
	mu       sync.RWMutex
}

// NewRegistry instantiates a new registry.
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]Service),
		list:     make([]string, 0),
	}
}

// Register a service using an id.
func (r *Registry) Register(id string, svc Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.services[id]
	if ok {
		return errors.New("service already exists")
	}
	r.services[id] = svc
	r.list = append(r.list, id)
	return nil
}

// GetService implements Discover interface.
func (r *Registry) GetService(id string) (Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.services[id]
	if !ok {
		return nil, false
	}
	return svc, true
}

// ListServices implements Discover interface.
func (r *Registry) ListServices() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]string, len(r.list), len(r.list))
	copy(list, r.list)
	return list
}

// Ping all registered services.
func (r *Registry) Ping() error {
	errs := make([]string, 0, len(r.services))
	for _, id := range r.ListServices() {
		svc, ok := r.GetService(id)
		if ok {
			err := svc.Ping()
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %s", id, err.Error()))
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}

// CloseAll registered services.
func (r *Registry) CloseAll() error {
	errs := make([]string, 0, len(r.services))
	for _, id := range r.ListServices() {
		svc, ok := r.GetService(id)
		if ok {
			err := svc.Close()
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
