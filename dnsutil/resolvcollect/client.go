// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package resolvcollect

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/dnsutil"
	pb "github.com/luids-io/core/protogen/dnsutilpb"
)

// Client implements a grpc client that implements dnsutil.ResolvCollector
// interface.
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.ResolvCollectClient
	//control
	started bool
}

type clientOpts struct {
	logger    yalogi.Logger
	closeConn bool
}

var defaultClientOpts = clientOpts{
	logger:    yalogi.LogNull,
	closeConn: true,
}

// ClientOption encapsules options for client
type ClientOption func(*clientOpts)

// CloseConnection option closes grpc connection on shutdown
func CloseConnection(b bool) ClientOption {
	return func(o *clientOpts) {
		o.closeConn = b
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

// NewClient returns a new client
func NewClient(conn *grpc.ClientConn, opt ...ClientOption) *Client {
	opts := defaultClientOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Client{
		opts:    opts,
		logger:  opts.logger,
		conn:    conn,
		client:  pb.NewResolvCollectClient(conn),
		started: true,
	}
}

// Collect implements dnsutil.ResolvCollector interface
func (c *Client) Collect(ctx context.Context, client net.IP, name string, resolved []net.IP) error {
	if !c.started {
		return errors.New("client closed")
	}
	return c.doCollect(ctx, client, name, resolved)
}

// Collect implements dnsutil.ResolvCollector interface
func (c *Client) doCollect(ctx context.Context, client net.IP, name string, resolved []net.IP) error {
	rr := make([]string, 0, len(resolved))
	for _, r := range resolved {
		rr = append(rr, r.String())
	}
	req := &pb.ResolvCollectRequest{
		ClientIp:    client.String(),
		Name:        name,
		ResolvedIps: rr,
	}
	_, err := c.client.Collect(ctx, req)
	if err != nil {
		return c.mapError(err)
	}
	return nil
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
		retErr = dnsutil.ErrBadRequestFormat
	case codes.ResourceExhausted:
		if st.Message() == dnsutil.ErrCollectDNSClientLimit.Error() {
			retErr = dnsutil.ErrCollectDNSClientLimit
		}
		if st.Message() == dnsutil.ErrCollectNamesLimit.Error() {
			retErr = dnsutil.ErrCollectNamesLimit
		}
	case codes.Unavailable:
		if st.Message() == dnsutil.ErrCacheNotAvailable.Error() {
			retErr = dnsutil.ErrCacheNotAvailable
		}
	}
	return retErr
}

//Close closes the client
func (c *Client) Close() error {
	if !c.started {
		return errors.New("client closed")
	}
	c.started = false
	if c.opts.closeConn {
		return c.conn.Close()
	}
	return nil
}

// Ping checks connectivity with the api
func (c *Client) Ping() error {
	if !c.started {
		return errors.New("client closed")
	}
	st := c.conn.GetState()
	switch st {
	case connectivity.TransientFailure:
		return fmt.Errorf("connection state: %v", st)
	case connectivity.Shutdown:
		return fmt.Errorf("connection state: %v", st)
	}
	return nil
}

//API returns API service name implemented
func (c *Client) API() string {
	return ServiceName()
}
