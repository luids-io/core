// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package serverd provides a component to create server-type applications,
// helping to manage the life cycle of services and operating system signals.
//
// This package is a work in progress and makes no API stability promises.
package serverd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/luids-io/core/yalogi"
)

// Option is used for component configuration.
type Option func(*options)

type options struct {
	logger          yalogi.Logger
	shutdownTimeout time.Duration
}

var defaultOptions = options{
	logger:          yalogi.LogNull,
	shutdownTimeout: 2 * time.Second,
}

// SetLogger option allows set a custom logger.
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// ShutdownTimeout option sets timeout for shutdowns.
func ShutdownTimeout(d time.Duration) Option {
	return func(o *options) {
		o.shutdownTimeout = d
	}
}

// Manager manages the services.
type Manager struct {
	opts   options
	logger yalogi.Logger

	name      string
	reloadErr error
	mu        sync.Mutex
	started   bool
	services  []Service
}

// New creates a new manager.
func New(name string, opt ...Option) *Manager {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	m := &Manager{
		opts:     opts,
		logger:   opts.logger,
		name:     name,
		services: make([]Service, 0),
	}
	return m
}

// Register resgisters a new service in the manager.
// Services cannot be registered if it has already been started.
func (m *Manager) Register(svc Service) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return errors.New("serverd: manager is started")
	}
	if svc.Name == "" {
		return errors.New("serverd: service name can't be empty")
	}
	for _, s := range m.services {
		if s.Name == svc.Name {
			return errors.New("serverd: service name can't be duplicated")
		}
	}
	m.services = append(m.services, svc)
	return nil
}

// Start starts managed services in the same order they were registered.
// If there is an error then it will stop.
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return nil
	}
	m.logger.Infof("starting %s services", m.name)
	for _, s := range m.services {
		m.logger.Infof("starting %s", s.Name)
		err := s.start()
		if err != nil {
			return fmt.Errorf("serverd: starting %s: %v", s.Name, err)
		}
	}
	m.started = true
	return nil
}

// Shutdown will stop all services in the opposite order to which they
// were registered. It will try to execute the "shutdown" function and,
// if it is turned off in a given time, it will execute the "stop" function
// of the service.
func (m *Manager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.started {
		return
	}
	m.started = false
	m.logger.Infof("shutting down %s services", m.name)
	for i := len(m.services) - 1; i >= 0; i-- {
		s := m.services[i]
		m.logger.Infof("shutting down %s", s.Name)
		s.shutdown(m.opts.shutdownTimeout)
	}
}

// Run will initialize all services, install the operating system's
// signal handlers and block waiting for the shutdown signal.
// When this signal arrives, it will turn off all registered services.
func (m *Manager) Run() error {
	// start services
	err := m.Start()
	if err != nil {
		return err
	}
	//launch signal handling goroutine
	close := make(chan bool, 1)
	go m.signalHndl(close)
	//waits for signal
	<-close
	// shutdown services
	m.Shutdown()
	return nil
}

func (m *Manager) signalHndl(close chan bool) {
	sigint := make(chan os.Signal, 1)
	sigHUP := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	signal.Notify(sigHUP, syscall.SIGHUP)
	for {
		select {
		case <-sigint:
			close <- true
			return
		case <-sigHUP:
			m.Reload()
		}
	}

}

// Ping will ping all registered services.
func (m *Manager) Ping() error {
	if !m.started {
		return errors.New("serverd: manager not started")
	}
	err := m.ping()
	if err != nil {
		return err
	}
	if m.reloadErr != nil {
		return m.reloadErr
	}
	return nil
}

// Reload will reload all registered services.
func (m *Manager) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.started {
		return nil
	}
	m.logger.Infof("reloading %s services", m.name)
	for _, s := range m.services {
		m.logger.Debugf("reloading service %s", s.Name)
		err := s.reload()
		if err != nil {
			m.logger.Warnf("reloading %s: %v", s.Name, err)
			m.reloadErr = fmt.Errorf("serverd: reloading %s: %v", s.Name, err)
			return m.reloadErr
		}
	}
	m.reloadErr = nil
	return nil
}

func (m *Manager) ping() error {
	errs := make([]string, 0, len(m.services))
	for _, s := range m.services {
		m.logger.Debugf("ping service %s", s.Name)
		err := s.ping()
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %s", s.Name, err.Error()))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}
