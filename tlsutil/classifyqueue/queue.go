// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package classifyqueue

import (
	"errors"
	"sync"
	"time"

	"github.com/luids-io/core/tlsutil"
)

type processQueueFn func(*sync.WaitGroup, []*tlsutil.ConnectionData)

type queue struct {
	wg       *sync.WaitGroup
	closed   bool
	reqC     chan *tlsutil.ConnectionData
	size     int
	interval time.Duration
	process  processQueueFn
}

func newQueue(s int, d time.Duration, p processQueueFn) *queue {
	q := &queue{
		wg:       &sync.WaitGroup{},
		size:     s,
		interval: d,
		reqC:     make(chan *tlsutil.ConnectionData, s),
		process:  p,
	}
	go q.run()
	return q
}

// close queue
func (q *queue) close() error {
	if q.closed {
		return errors.New("queue already closed")
	}
	q.closed = true
	close(q.reqC)
	q.wg.Wait()
	return nil
}

// add to queue
func (q *queue) add(req *tlsutil.ConnectionData) error {
	if q.closed {
		return errors.New("queue is closed")
	}
	q.reqC <- req
	return nil
}

// run is the main loop
func (q *queue) run() {
	buffer := make([]*tlsutil.ConnectionData, 0, q.size)
	ticker := time.NewTicker(q.interval)
	defer ticker.Stop()
	for {
		select {
		case request, ok := <-q.reqC:
			if !ok {
				return
			}
			buffer = append(buffer, request)
			if len(buffer) >= q.size {
				//process job in goroutine
				q.wg.Add(1)
				go q.process(q.wg, buffer)
				buffer = make([]*tlsutil.ConnectionData, 0, q.size)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				//process job in goroutine
				q.wg.Add(1)
				go q.process(q.wg, buffer)
				buffer = make([]*tlsutil.ConnectionData, 0, q.size)
			}
		}
	}
}
