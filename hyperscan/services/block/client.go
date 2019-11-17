// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package block

import (
	"context"
	"errors"
	"fmt"

	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	pb "github.com/luids-io/core/protogen/hspb"
)

// Client provides a grpc client that implements hyperscan.BlockScanner interface.
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.BlockClient
	//control
	started bool
}

type clientOpts struct {
	logger    yalogi.Logger
	closeConn bool
	debugreq  bool
}

var defaultClientOpts = clientOpts{
	logger:    yalogi.LogNull,
	closeConn: true,
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

// NewClient returns a new grpc Client
func NewClient(conn *grpc.ClientConn, opt ...ClientOption) *Client {
	opts := defaultClientOpts
	for _, o := range opt {
		o(&opts)
	}
	c := &Client{
		opts:   opts,
		logger: opts.logger,
		conn:   conn,
		client: pb.NewBlockClient(conn),
	}
	c.started = true
	return c
}

// ScanBlock implements hyperscan.BlockScanner interface
func (c *Client) ScanBlock(ctx context.Context, block []byte) (bool, []string, error) {
	if !c.started {
		return false, nil, errors.New("client not started")
	}
	if c.opts.debugreq {
		c.logger.Debugf("scan(%v)", block)
	}
	return c.doScan(ctx, block)
}

func (c *Client) doScan(ctx context.Context, data []byte) (bool, []string, error) {
	req := &pb.ScanBlockRequest{Data: data}
	res, err := c.client.ScanBlock(ctx, req)
	if err != nil {
		return false, nil, c.mapError(err)
	}
	return res.GetResult(), res.GetReasons(), nil
}

//mapping errors
func (c *Client) mapError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	return errors.New(st.Message())
}

//Close the client
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
