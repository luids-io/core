// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/dnsutil"
	pb "github.com/luids-io/core/protogen/dnsutilpb"
)

// Client implements a grpc client that implements dnsutil.ResolvArchiver
// interface.
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.ArchiveClient
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
		client:  pb.NewArchiveClient(conn),
		started: true,
	}
}

// SaveResolv implements dnsutil.Archiver interface
func (c *Client) SaveResolv(ctx context.Context, data dnsutil.ResolvData) (string, error) {
	if !c.started {
		return "", errors.New("client closed")
	}
	return c.doSave(ctx, data)
}

func (c *Client) doSave(ctx context.Context, data dnsutil.ResolvData) (string, error) {
	tstamp, _ := ptypes.TimestampProto(data.Timestamp)
	rr := make([]string, 0, len(data.Resolved))
	for _, r := range data.Resolved {
		rr = append(rr, r.String())
	}
	req := &pb.SaveResolvRequest{
		Ts:          tstamp,
		ServerIp:    data.Server.String(),
		ClientIp:    data.Client.String(),
		Name:        data.Name,
		ResolvedIps: rr,
	}
	resp, err := c.client.SaveResolv(ctx, req)
	if err != nil {
		return "", c.mapError(err)
	}
	return resp.GetId(), nil
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
	case codes.Unavailable:
		if st.Message() == dnsutil.ErrServiceNotAvailable.Error() {
			retErr = dnsutil.ErrServiceNotAvailable
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
