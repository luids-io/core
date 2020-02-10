// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package manager

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/core/brain/classify"
)

type options struct {
	logger   yalogi.Logger
	queueLen int
	maxItems int
	interval time.Duration
}

var defaultOpts = options{
	logger:   yalogi.LogNull,
	queueLen: 1024,
	maxItems: 100,
	interval: 1 * time.Second,
}

// Option encapsules options for manager
type Option func(*options)

// SetLogger option allows set a custom logger
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// Dispatcher function for classification responses
type Dispatcher func([]classify.Request, []classify.Response)

// Manager implements a manager for a classification service
type Manager struct {
	opts       options
	logger     yalogi.Logger
	classifier classify.Classifier
	dispatch   Dispatcher
	wg         sync.WaitGroup
	closed     bool
	reqC       chan classify.Request
}

// New returns a new classify manager
func New(c classify.Classifier, opt ...Option) *Manager {
	opts := defaultOpts
	for _, o := range opt {
		o(&opts)
	}
	m := &Manager{
		opts:       opts,
		logger:     opts.logger,
		classifier: c,
		reqC:       make(chan classify.Request, opts.queueLen),
	}
	go m.run()
	return m
}

// SetDispatcher to dispatch requests
func (m *Manager) SetDispatcher(d Dispatcher) {
	m.dispatch = d
}

// Push a new request in the classification queue
func (m *Manager) Push(req classify.Request) error {
	if m.closed {
		return errors.New("classify manager is closed")
	}
	m.reqC <- req
	return nil
}

// Close manager
func (m *Manager) Close() error {
	if m.closed {
		return errors.New("classify manager already closed")
	}
	m.closed = true
	close(m.reqC)
	return nil
}

// run is the main loop
func (m *Manager) run() {
	queue := make([]classify.Request, 0, m.opts.maxItems)
	ticker := time.NewTicker(m.opts.interval)
	defer ticker.Stop()
	for {
		select {
		case request, ok := <-m.reqC:
			if !ok {
				return
			}
			queue = append(queue, request)
			if len(queue) >= m.opts.maxItems {
				//process job in goroutine
				m.wg.Add(1)
				go m.process(queue)
				queue = make([]classify.Request, 0, m.opts.maxItems)
			}
		case <-ticker.C:
			if len(queue) > 0 {
				//process job in goroutine
				m.wg.Add(1)
				go m.process(queue)
				queue = make([]classify.Request, 0, m.opts.maxItems)
			}
		}
	}
}

func (m *Manager) process(requests []classify.Request) {
	m.logger.Debugf("classifying requests: len(requests)=%v", len(requests))
	defer m.wg.Done()
	//TODO: implementar cancelación
	responses, err := m.classifier.Classify(context.Background(), requests)
	if err != nil {
		m.logger.Warnf("classifying requests: %v", err)
		return
	}
	if m.dispatch != nil {
		m.logger.Debugf("dispatching requests len(requests)=%v", len(requests))
		m.dispatch(requests, responses)
	}
}
