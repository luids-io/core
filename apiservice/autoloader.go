// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package apiservice

import (
	"sort"
	"sync"

	"github.com/luids-io/core/yalogi"
)

// Autoloader implements Discover using a lazy build.
type Autoloader struct {
	logger yalogi.Logger
	defs   map[string]ServiceDef
	reg    *Registry
	mu     sync.RWMutex
}

// AutoloaderOption is used for Autoloader configuration.
type AutoloaderOption func(*autoOptions)

type autoOptions struct {
	logger yalogi.Logger
}

var defaultAutoOptions = autoOptions{logger: yalogi.LogNull}

// SetLogger option allows set a custom logger.
func SetLogger(l yalogi.Logger) AutoloaderOption {
	return func(o *autoOptions) {
		if l != nil {
			o.logger = l
		}
	}
}

// NewAutoloader creates a new Autoloader with service definitions.
func NewAutoloader(defs []ServiceDef, opt ...AutoloaderOption) *Autoloader {
	opts := defaultAutoOptions
	for _, o := range opt {
		o(&opts)
	}
	a := &Autoloader{
		logger: opts.logger,
		defs:   make(map[string]ServiceDef),
		reg:    NewRegistry(),
	}
	for _, def := range defs {
		if !def.Disabled {
			a.defs[def.ID] = def
		}
	}
	return a
}

// GetService implements Discover interface.
func (a *Autoloader) GetService(id string) (Service, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	//try get from registry
	svc, ok := a.reg.GetService(id)
	if ok {
		return svc, true
	}
	//get definition
	def, ok := a.defs[id]
	if !ok {
		return nil, false
	}
	//build service & register
	svc, err := Build(def, a.logger)
	if err != nil {
		a.logger.Errorf("apiservice: autoloader building service '%s': %v", id, err)
		return nil, false
	}
	a.reg.Register(id, svc)
	return svc, ok
}

// ListServices implements Discover interface.
func (a *Autoloader) ListServices() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	list := make([]string, 0, len(a.defs))
	for k := range a.defs {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// Ping all registered services.
func (a *Autoloader) Ping() error {
	return a.reg.Ping()
}

// CloseAll registered services.
func (a *Autoloader) CloseAll() error {
	return a.reg.CloseAll()
}
