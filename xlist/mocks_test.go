// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist_test

import (
	"context"
	"fmt"

	"github.com/luids-io/core/xlist"
)

type mockList struct {
	fail      bool
	resources []xlist.Resource
	response  xlist.Response
	//lazy      time.Duration
}

func (l mockList) Check(ctx context.Context, name string, res xlist.Resource) (xlist.Response, error) {
	name, ctx, err := xlist.DoValidation(ctx, name, res, false)
	if err != nil {
		return xlist.Response{}, err
	}
	if !res.InArray(l.resources) {
		return xlist.Response{}, xlist.ErrResourceNotSupported
	}
	if l.fail {
		return xlist.Response{}, xlist.ErrListNotAvailable
	}
	// if l.lazy > 0 {
	// 	time.Sleep(l.lazy)
	// }
	return l.response, nil
}

func (l mockList) Ping() error {
	if l.fail {
		return xlist.ErrListNotAvailable
	}
	return nil
}

func (l mockList) Resources() []xlist.Resource {
	ret := make([]xlist.Resource, len(l.resources), len(l.resources))
	copy(ret, l.resources)
	return ret
}

type mockContainer struct {
	stopOnError bool
	resources   []xlist.Resource
	lists       []xlist.Checker
}

func (c mockContainer) Check(ctx context.Context, name string, res xlist.Resource) (xlist.Response, error) {
	name, ctx, err := xlist.DoValidation(ctx, name, res, false)
	if err != nil {
		return xlist.Response{}, err
	}
	if !res.InArray(c.resources) {
		return xlist.Response{}, xlist.ErrResourceNotSupported
	}
	for _, checker := range c.lists {
		resp, err := checker.Check(ctx, name, res)
		if err != nil && c.stopOnError {
			return xlist.Response{}, err
		}
		if resp.Result {
			return resp, err
		}
	}
	return xlist.Response{}, nil
}

func (c mockContainer) Ping() error {
	for _, checker := range c.lists {
		err := checker.Ping()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c mockContainer) Resources() []xlist.Resource {
	var ret []xlist.Resource
	copy(ret, c.resources)
	return ret
}

type mockWrapper struct {
	preffix string
	checker xlist.Checker
}

func (w mockWrapper) Check(ctx context.Context, name string, res xlist.Resource) (xlist.Response, error) {
	resp, err := w.checker.Check(ctx, name, res)
	resp.Reason = fmt.Sprintf("%s: %s", w.preffix, resp.Reason)
	return resp, err
}

func (w mockWrapper) Ping() error {
	return w.checker.Ping()
}

func (w mockWrapper) Resources() []xlist.Resource {
	return w.Resources()
}
