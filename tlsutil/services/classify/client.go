// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package classify

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/luids-io/core/brain/classify"
	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
	"github.com/luids-io/core/tlsutil/services/encoding"
	"github.com/luisguillenc/yalogi"
)

// Client provides a grpc client that implements a tlsutil machine learning classifier
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.ClassifyClient
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
		client: pb.NewClassifyClient(conn),
	}
	c.started = true
	return c
}

// Classify implements classify.Classifier
func (c *Client) Classify(ctx context.Context, requests []classify.Request) ([]classify.Response, error) {
	if !c.started {
		return nil, errors.New("client closed")
	}
	// prepare requests
	sendRequests := make([]*pb.ClassifyConnectionRequest_Request, 0, len(requests))
	for _, r := range requests {
		cdata, ok := r.Data.(*tlsutil.ConnectionData)
		if !ok {
			return nil, fmt.Errorf("request id %s is not valid", r.ID)
		}
		sendRequests = append(sendRequests, &pb.ClassifyConnectionRequest_Request{
			Id:         r.ID,
			Connection: encoding.ConnectionDataPB(cdata),
		})
	}
	// do classify
	pbres, err := c.client.Connections(ctx, &pb.ClassifyConnectionRequest{Requests: sendRequests})
	if err != nil {
		return nil, err
	}
	if len(requests) != len(pbres.Responses) {
		return nil, errors.New("requests len and responses len missmatch")
	}
	// reencode responses
	responses := make([]classify.Response, 0, len(pbres.Responses))
	for _, r := range pbres.Responses {
		resp := classify.Response{}
		resp.ID = r.GetId()
		if r.GetErr() != "" {
			resp.Err = errors.New(r.GetErr())
		} else {
			resp.Results = make([]classify.Result, 0, len(r.GetResults()))
			for _, result := range r.GetResults() {
				resp.Results = append(resp.Results, classify.Result{
					Label: result.GetLabel(),
					Prob:  result.GetProb(),
				})
			}
		}
		responses = append(responses, resp)
	}
	return responses, nil
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
