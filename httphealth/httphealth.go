// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.md.

// Package httphealth provides a simple component that offers an http interface
// for health checking and monitoring using Prometheus.
//
// This package is a work in progress and makes no API stability promises.
package httphealth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/luids-io/core/ipfilter"
	"github.com/luids-io/core/yalogi"
)

// Pingable must be implemented by the service to be monitored
type Pingable interface {
	Ping() error
}

type options struct {
	logger   yalogi.Logger
	ipfilter ipfilter.Filter
	metrics  bool
	profile  bool
}

var defaultOptions = options{logger: yalogi.LogNull}

// Option encapsules options for server
type Option func(*options)

// SetLogger sets a logger for the component
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// SetIPFilter sets an ip filter for the health server
func SetIPFilter(f ipfilter.Filter) Option {
	return func(o *options) {
		o.ipfilter = f
	}
}

// Metrics enabled to expose prometheus metrics in the health server
func Metrics(b bool) Option {
	return func(o *options) {
		o.metrics = b
	}
}

// Profile enabled to expose pprof in the health server
func Profile(b bool) Option {
	return func(o *options) {
		o.profile = b
	}
}

// Server is an http server wrapped that provides a health service
type Server struct {
	opts       options
	logger     yalogi.Logger
	server     *http.Server
	supervised Pingable
}

// New construct a new health server that supervised the 'Pingable' object
func New(supervised Pingable, opt ...Option) *Server {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &Server{
		opts:       opts,
		logger:     opts.logger,
		server:     &http.Server{},
		supervised: supervised,
	}
	return s
}

// Serve calls the same function from the http.Server contained in the struct
func (s *Server) Serve(lis net.Listener) error {
	s.logger.Infof("starting health server %v", lis.Addr().String())
	s.server.Handler = s.handler()
	return s.server.Serve(lis)
}

// ServeTLS calls the same function from the http.Server contained in the struct
func (s *Server) ServeTLS(lis net.Listener, certFile string, keyFile string) error {
	s.logger.Infof("starting health server (tls) %v", lis.Addr().String())
	s.server.Handler = s.handler()
	return s.server.ServeTLS(lis, certFile, keyFile)
}

// Close calls the same function from the http.Server contained in the struct
func (s *Server) Close() error {
	s.logger.Infof("closing health server")
	return s.server.Close()
}

// Shutdown calls the same function from the http.Server contained in the struct
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Infof("shutting down health server")
	return s.server.Shutdown(ctx)
}

func (s *Server) handler() http.Handler {
	router := mux.NewRouter()
	if s.opts.metrics {
		router.Handle("/metrics", promhttp.Handler())
	}
	if s.opts.profile {
		attachProfiler(router)
	}
	router.HandleFunc("/health", s.doHealth).Methods("GET")

	if !s.opts.ipfilter.Empty() {
		filtered := s.opts.ipfilter
		filtered.Wrapped = router
		return filtered
	}
	return router
}

func (s *Server) doHealth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := s.supervised.Ping()
	latency := time.Since(start)
	if err != nil {
		fmt.Fprintf(w, "status: FAILED (%s)\n", err.Error())
		s.logger.Debugf("health request from %s return error: %v", r.RemoteAddr, err)
	} else {
		fmt.Fprintf(w, "status: OK\n")
		s.logger.Debugf("health request from %s return ok", r.RemoteAddr)
	}
	fmt.Fprintf(w, "latency: %v", latency)
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
