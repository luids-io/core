// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package classifyqueue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/luids-io/core/tlsutil"
	"github.com/luids-io/core/utils/yalogi"
)

// Dispatcher function for process classification responses
type Dispatcher func([]*tlsutil.ConnectionData, []tlsutil.ClassifyResponse)

// Manager implements an async classification service
type Manager struct {
	logger yalogi.Logger
	closed bool
	//tls calssifier
	classifier tlsutil.Classifier
	//queues for classification
	queueCon *queue
	//dispatchers
	dispCon Dispatcher
}

type options struct {
	logger    yalogi.Logger
	queueSize int
	interval  time.Duration
}

var defaultOpts = options{
	logger:    yalogi.LogNull,
	queueSize: 128,
	interval:  1 * time.Second,
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

// SetQueueSize option allows change classification queue size
func SetQueueSize(i int) Option {
	return func(o *options) {
		if i > 0 {
			o.queueSize = i
		}
	}
}

// SetInterval option allows change interval
func SetInterval(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.interval = d
		}
	}
}

// New returns a new classify manager
func New(c tlsutil.Classifier, opt ...Option) *Manager {
	opts := defaultOpts
	for _, o := range opt {
		o(&opts)
	}
	m := &Manager{
		logger:     opts.logger,
		classifier: c,
	}
	//create queues
	m.queueCon = newQueue(opts.queueSize, opts.interval, m.classifyConnections)
	return m
}

// SetConnectionDispatcher to dispatch connection classify responses
func (m *Manager) SetConnectionDispatcher(d Dispatcher) {
	m.dispCon = d
}

// PushConnection add a new connection in the classification queue
func (m *Manager) PushConnection(req *tlsutil.ConnectionData) error {
	if m.closed {
		return errors.New("manager is closed")
	}
	return m.queueCon.add(req)
}

// Close manager
func (m *Manager) Close() error {
	if m.closed {
		return errors.New("manager already closed")
	}
	m.closed = true
	m.queueCon.close()
	return nil
}

// classifyConnections implements a processQueueFn function
func (m *Manager) classifyConnections(wg *sync.WaitGroup, requests []*tlsutil.ConnectionData) {
	m.logger.Debugf("classifying connections: len(requests)=%v", len(requests))
	defer wg.Done()
	//TODO: implement cancel
	responses, err := m.classifier.ClassifyConnections(context.Background(), requests)
	if err != nil {
		m.logger.Warnf("classifying connections: %v", err)
		return
	}
	if m.dispCon != nil {
		m.logger.Debugf("dispatching responses len(responses)=%v", len(responses))
		m.dispCon(requests, responses)
	}
}
