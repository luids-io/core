// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package packetproc

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/gopacket"

	"github.com/luids-io/core/utils/yalogi"
)

// Processor defines main Process function
type Processor interface {
	Process(name string, source gopacket.PacketDataSource, hooks *Hooks) (stop func(), errs <-chan error, err error)
}

// Service defines a packet processor service
type Service interface {
	Start() error
	Shutdown()
	Ping() error
	AddPlugin(p Plugin) error
	Register(name string, source gopacket.PacketDataSource, proc Processor) error
	Unregister(name string) error
}

// service packet processor implementation
type service struct {
	sources map[string]*pcktSource
	names   map[string]bool
	plugins []Plugin
	logger  yalogi.Logger
	//control
	mu      sync.Mutex
	wg      sync.WaitGroup
	started bool
	errCh   chan error
}

type options struct {
	logger yalogi.Logger
}

var defaultOptions = options{
	logger: yalogi.LogNull,
}

// Option is used for component configuration
type Option func(*options)

// SetLogger option allows set a custom logger
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// NewService creates a new service
func NewService(opt ...Option) Service {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &service{
		logger:  opts.logger,
		sources: make(map[string]*pcktSource),
		names:   make(map[string]bool),
		plugins: make([]Plugin, 0),
	}
	return s
}

// Start the service and start to process registered packet sources
func (s *service) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return errors.New("service started")
	}
	s.logger.Infof("starting packet processing service")
	// create errors channel and process it
	s.errCh = make(chan error, ErrorsBuffer)
	go s.procErrs()
	s.started = true
	// start processing all registered sources
	for _, src := range s.sources {
		err := s.doStart(src)
		if err != nil {
			return err
		}
	}
	return nil
}

// Shutdown the service and stop processing registered packet sources
func (s *service) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started {
		return
	}
	s.logger.Infof("shutting down packet processing service")
	for _, src := range s.sources {
		if src.started {
			src.stop()
		}
	}
	s.wg.Wait()
	close(s.errCh)
	s.started = false
}

type pcktSource struct {
	name    string
	source  gopacket.PacketDataSource
	proc    Processor
	started bool
	stop    func()
	errCh   <-chan error
}

// AddPlugin to service
func (s *service) AddPlugin(p Plugin) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pname := p.Name()
	_, repeated := s.names[pname]
	if repeated {
		return fmt.Errorf("plugin '%s' already registered", pname)
	}
	s.logger.Debugf("registering plugin '%s' class '%s'", p.Name(), p.Class())
	s.plugins = append(s.plugins, p)
	s.names[pname] = true
	return nil
}

// Register packet source with name and start it if service is started
func (s *service) Register(name string, source gopacket.PacketDataSource, proc Processor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debugf("registering source %s", name)
	_, ok := s.sources[name]
	if ok {
		return errors.New("packet source exists")
	}
	src := &pcktSource{
		name:   name,
		source: source,
		proc:   proc,
	}
	s.sources[name] = src
	if s.started {
		return s.doStart(src)
	}
	return nil
}

// Unregister packet source by name, stopping if it's started
func (s *service) Unregister(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debugf("unregistering packet source %s", name)
	src, ok := s.sources[name]
	if !ok {
		return errors.New("packet source doesn't exists")
	}
	if s.started && src.started {
		src.stop()
	}
	delete(s.sources, name)
	return nil
}

// Ping returns true if errors
func (s *service) Ping() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started {
		return errors.New("service not started")
	}
	errs := make([]string, 0, len(s.sources))
	for _, src := range s.sources {
		if !src.started {
			errs = append(errs, src.name)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("packet sources stopped: %s", strings.Join(errs, ","))
	}
	return nil
}

//start packet source
func (s *service) doStart(src *pcktSource) error {
	s.logger.Infof("starting packet source %s", src.name)
	hooks := NewHooks()
	for _, p := range s.plugins {
		p.Register(src.name, hooks)
	}
	var err error
	src.stop, src.errCh, err = src.proc.Process(src.name, src.source, hooks)
	if err != nil {
		return fmt.Errorf("starting packet source %s: %v", src.name, err)
	}
	src.started = true
	s.wg.Add(1)
	//processing error channel goroutine
	go func(c <-chan error) {
		for n := range c {
			s.errCh <- n
		}
		s.logger.Infof("stopping packet source %s", src.name)
		src.started = false
		s.wg.Done()
	}(src.errCh)
	return nil
}

//routine for processing error services
func (s *service) procErrs() {
	for err := range s.errCh {
		s.logger.Warnf("%v", err)
	}
}
