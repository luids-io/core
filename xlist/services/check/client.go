// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package check

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/luids-io/core/protogen/xlistpb"
	"github.com/luids-io/core/xlist"
)

// Client provides a grpc client that implements xlist.Checker interface.
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.CheckClient
	//control
	started, synced bool
	mu, musync      sync.Mutex
	resources       []xlist.Resource
	//cache
	cache *cache
}

type clientOpts struct {
	logger          yalogi.Logger
	closeConn       bool
	forceValidation bool
	debugreq        bool
	//cache opts
	useCache     bool
	ttl          int
	negativettl  int
	cacheCleanup time.Duration
}

var defaultClientOpts = clientOpts{
	logger:       yalogi.LogNull,
	closeConn:    true,
	cacheCleanup: defaultCacheCleanups,
}

// ClientOption encapsules options for client
type ClientOption func(*clientOpts)

// CloseConnection option closes grpc connection on close
func CloseConnection(b bool) ClientOption {
	return func(o *clientOpts) {
		o.closeConn = b
	}
}

// DebugRequests option enables debug messages in requests
func DebugRequests(b bool) ClientOption {
	return func(o *clientOpts) {
		o.debugreq = b
	}
}

// SetLogger option allows set a custom logger
func SetLogger(l yalogi.Logger) ClientOption {
	return func(o *clientOpts) {
		if l != nil {
			o.logger = l
		}
	}
}

// ForceValidation forces component to ignore context and validate requests
func ForceValidation(b bool) ClientOption {
	return func(o *clientOpts) {
		o.forceValidation = b
	}
}

// SetCache sets cache ttl and negative ttl
func SetCache(ttl, negativettl int) ClientOption {
	return func(o *clientOpts) {
		if ttl >= xlist.NeverCache && negativettl >= xlist.NeverCache {
			o.ttl = ttl
			o.negativettl = negativettl
			o.useCache = true
		}
	}
}

// SetCacheCleanUps sets interval between cache cleanups
func SetCacheCleanUps(d time.Duration) ClientOption {
	return func(o *clientOpts) {
		if d > 0 {
			o.cacheCleanup = d
		}
	}
}

// NewClient returns a new grpc Client
func NewClient(conn *grpc.ClientConn, resources []xlist.Resource, opt ...ClientOption) *Client {
	opts := defaultClientOpts
	for _, o := range opt {
		o(&opts)
	}
	c := &Client{
		opts:   opts,
		logger: opts.logger,
		conn:   conn,
		client: pb.NewCheckClient(conn),
	}
	if len(resources) > 0 {
		c.resources = xlist.ClearResourceDups(resources)
		c.synced = true
	}
	//if no resources passed, it will get the resources supported
	//by the checker in the first check -> synced = false
	if opts.useCache {
		c.cache = newCache(opts.ttl, opts.negativettl, opts.cacheCleanup)
	}
	c.started = true
	return c
}

// Check implements xlist.Checker interface
func (c *Client) Check(ctx context.Context, name string, resource xlist.Resource) (xlist.Response, error) {
	if !c.started {
		return xlist.Response{}, xlist.ErrListNotAvailable
	}
	if c.opts.debugreq {
		c.logger.Debugf("check(%s,%v)", name, resource)
	}
	if !c.synced {
		if err := c.sync(ctx); err != nil {
			return xlist.Response{}, xlist.ErrListNotAvailable
		}
	}
	if !c.checks(resource) {
		return xlist.Response{}, xlist.ErrResourceNotSupported
	}
	name, ctx, err := xlist.DoValidation(ctx, name, resource, c.opts.forceValidation)
	if err != nil {
		return xlist.Response{}, err
	}
	if c.opts.useCache {
		resp, ok := c.cache.get(name, resource)
		if ok {
			return resp, nil
		}
	}
	resp, err := c.doCheck(ctx, name, resource)
	if c.opts.useCache {
		if err == nil {
			c.cache.set(name, resource, resp)
		}
	}
	return resp, err
}

// Resources implements xlist.Checker interface
func (c *Client) Resources() []xlist.Resource {
	if c.opts.debugreq {
		c.logger.Debugf("resources()")
	}
	if !c.synced && c.started {
		c.sync(context.Background())
	}
	return c.resources
}

// Ping implements xlist.Checker interface
func (c *Client) Ping() error {
	if !c.started {
		return errors.New("client not started")
	}
	if c.opts.debugreq {
		c.logger.Debugf("ping()")
	}
	return c.doPing(context.Background())
}

//sync resources
func (c *Client) sync(ctx context.Context) error {
	err := c.doPing(ctx)
	if err == nil {
		c.musync.Lock()
		defer c.musync.Unlock()
		if !c.synced {
			res := c.doResources(ctx)
			c.resources = res
			c.synced = true
			c.logger.Debugf("resources synced: %v", res)
		}
	}
	return err
}

func (c *Client) doCheck(ctx context.Context, name string, resource xlist.Resource) (xlist.Response, error) {
	req := &pb.CheckRequest{Request: &pb.Request{Name: name, Resource: pb.Resource(resource)}}
	res, err := c.client.Check(ctx, req)
	if err != nil {
		return xlist.Response{}, c.mapError(err)
	}
	r := xlist.Response{
		Result: res.GetResponse().GetResult(),
		Reason: res.GetResponse().GetReason(),
		TTL:    int(res.GetResponse().GetTTL())}
	return r, nil
}

func (c *Client) doResources(ctx context.Context) []xlist.Resource {
	resp, err := c.client.Resources(ctx, &empty.Empty{})
	if err != nil {
		return []xlist.Resource{}
	}
	resources := make([]xlist.Resource, 0, len(resp.Resources))
	for _, r := range resp.Resources {
		resources = append(resources, xlist.Resource(r))
	}
	return xlist.ClearResourceDups(resources)
}

func (c *Client) doPing(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &empty.Empty{})
	return err
}

//mapping errors
func (c *Client) mapError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	retErr := errors.New(st.Message())
	switch st.Code() {
	case codes.InvalidArgument:
		retErr = xlist.ErrBadResourceFormat
	case codes.Unimplemented:
		retErr = xlist.ErrResourceNotSupported
	case codes.Unavailable:
		if st.Message() == xlist.ErrListNotAvailable.Error() {
			retErr = xlist.ErrListNotAvailable
		}
	}
	return retErr
}

func (c *Client) checks(r xlist.Resource) bool {
	for _, res := range c.resources {
		if r == res {
			return true
		}
	}
	return false
}

//Flush cache if set
func (c *Client) Flush() {
	c.logger.Debugf("flushing cache")
	if c.opts.useCache {
		c.cache.flush()
	}
}

//Close the client
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.started {
		return errors.New("client closed")
	}
	c.logger.Debugf("closing connection")
	c.started = false
	if c.opts.closeConn {
		return c.conn.Close()
	}
	return nil
}

//API returns API service name implemented
func (c *Client) API() string {
	return ServiceName()
}
