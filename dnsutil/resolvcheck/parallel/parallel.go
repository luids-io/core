// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package parallel allows multiple checks in paralell using
// dnsutil.ResolvChecker interface.
//
// This package is a work in progress and makes no API stability promises.
package parallel

import (
	"context"
	"net"
	"sync"

	"github.com/luids-io/core/dnsutil"
)

// Request encapsules dnsutil resolv queries in a struct for parallel checks
type Request struct {
	Client   net.IP
	Resolved net.IP
	Name     string
}

// Response is used for store request and response in an Check parallel
// to the resolv checker
type Response struct {
	Request  Request
	Response dnsutil.ResolvResponse
	Err      error
}

// Check checks a set of queries in parallel. It returns an slice with the response,
// a boolean returning true if there is an error in at least one query and an error if
// there was a problem in the parallel check.
func Check(ctx context.Context, checkers []dnsutil.ResolvChecker,
	requests []Request) ([]Response, bool, error) {
	// prepare queries and launch workers
	var wg sync.WaitGroup
	ctxChild, cancelChilds := context.WithCancel(context.Background())
	results := make(chan Response, len(checkers)*len(requests))
	for _, checker := range checkers {
		for _, query := range requests {
			wg.Add(1)
			go workerCheck(ctxChild, &wg, checker,
				query.Client, query.Resolved, query.Name, results)
		}
	}
	// get results
	var err error
	returned := make([]Response, 0, len(requests))
	finished, hasErrors := 0, false
RESULTLOOP:
	for finished < len(requests) {
		select {
		case result := <-results:
			finished++
			if result.Err != nil {
				hasErrors = true
			}
			returned = append(returned, result)
		case <-ctx.Done():
			err = ctx.Err()
			break RESULTLOOP
		}
	}
	cancelChilds()
	wg.Wait()
	close(results)

	return returned, hasErrors, err
}

func workerCheck(ctx context.Context, wg *sync.WaitGroup, checker dnsutil.ResolvChecker,
	client, resolved net.IP, name string, results chan<- Response) {

	defer wg.Done()
	response, err := checker.Check(ctx, client, resolved, name)
	if err != nil {
		results <- Response{
			Request: Request{Client: client, Resolved: resolved, Name: name},
			Err:     err,
		}
		return
	}
	results <- Response{
		Request:  Request{Client: client, Resolved: resolved, Name: name},
		Response: response,
	}
	return
}
