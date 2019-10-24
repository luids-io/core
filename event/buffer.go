// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package event

import (
	"context"
	"errors"

	"github.com/luisguillenc/yalogi"
)

// Buffer implements a buffer for async event notification
type Buffer struct {
	//logger used for errors
	logger yalogi.Logger
	//collector
	notifier Notifier
	//data channel
	eventCh chan Event
	//control
	closed bool
	close  chan struct{}
}

type bufferOpts struct {
	logger yalogi.Logger
}

var defaultBufferOpts = bufferOpts{
	logger: yalogi.LogNull,
}

// BufferOption encapsules options for buffer
type BufferOption func(*bufferOpts)

// SetLogger option allows set a custom logger
func SetLogger(l yalogi.Logger) BufferOption {
	return func(o *bufferOpts) {
		if l != nil {
			o.logger = l
		}
	}
}

// NewBuffer returns a new event buffer
func NewBuffer(n Notifier, size int, opt ...BufferOption) *Buffer {
	opts := defaultBufferOpts
	for _, o := range opt {
		o(&opts)
	}
	b := &Buffer{
		logger:   opts.logger,
		notifier: n,
		eventCh:  make(chan Event, size),
		close:    make(chan struct{}),
	}
	go b.doProcess()
	return b
}

// Notify implements an asyncronous notification
func (b *Buffer) Notify(e Event) error {
	if b.closed {
		return errors.New("buffer is closed")
	}
	b.eventCh <- e
	return nil
}

func (b *Buffer) doProcess() {
	for e := range b.eventCh {
		reqid, err := b.notifier.Notify(context.Background(), e)
		if err != nil {
			b.logger.Warnf("%v", err)
		}
		b.logger.Debugf("notified event reqid: '%s'", reqid)
	}
	close(b.close)
}

// Close buffer
func (b *Buffer) Close() {
	if b.closed {
		return
	}
	b.closed = true
	close(b.eventCh)
	<-b.close
}
