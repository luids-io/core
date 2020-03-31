// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package serverd

import (
	"time"
)

// Service defines data struct for services
type Service struct {
	Name     string
	Start    StartupFn
	Shutdown ShutdownFn
	Stop     StopFn
	Reload   ReloadFn
	Ping     PingFn
}

// StartupFn sets definition for startup functions
type StartupFn func() error

// ShutdownFn sets definition for gracely shutdown functions
type ShutdownFn func()

// StopFn sets definition for stop functions (force)
type StopFn func()

// ReloadFn sets definition for reload functions
type ReloadFn func() error

// PingFn sets definition for ping functions
type PingFn func() error

func (s Service) start() error {
	var err error
	if s.Start != nil {
		err = s.Start()
	}
	return err
}

func (s Service) shutdown(timeout time.Duration) {
	if s.Shutdown != nil {
		fin := make(chan struct{})
		go func() {
			s.Shutdown()
			close(fin)
		}()
		select {
		case <-time.After(timeout):
			if s.Stop != nil {
				s.Stop()
			}
		case <-fin:
		}
	} else if s.Stop != nil {
		s.Stop()
	}
	return
}

func (s Service) ping() error {
	var err error
	if s.Ping != nil {
		err = s.Ping()
	}
	return err
}

func (s Service) reload() error {
	var err error
	if s.Reload != nil {
		err = s.Reload()
	}
	return err
}
