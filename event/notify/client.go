// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/luids-io/core/event"
	pb "github.com/luids-io/core/protogen/eventpb"
)

// Client is the main struct for grpc client
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.NotifyClient
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
		client:  pb.NewNotifyClient(conn),
		started: true,
	}
}

// Notify implements event.Notifier interface
func (c *Client) Notify(ctx context.Context, e event.Event) (string, error) {
	if !c.started {
		return "", errors.New("client closed")
	}
	return c.doNotify(ctx, e)
}

func (c *Client) doNotify(ctx context.Context, e event.Event) (string, error) {
	//create request
	req, err := eventToRequest(e)
	if err != nil {
		return "", fmt.Errorf("serializing event: %v", err)
	}
	//notify request
	resp, err := c.client.Notify(ctx, req)
	if err != nil {
		return "", c.mapError(err)
	}
	//process response
	requestID := resp.GetRequestID()
	if requestID == "" {
		return "", fmt.Errorf("processing response: request_id empty")
	}
	return requestID, nil
}

//mapping errors routine
func (c *Client) mapError(err error) error {
	//TODO
	return err
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

func eventToRequest(e event.Event) (*pb.NotifyRequest, error) {
	req := &pb.NotifyRequest{}
	req.Id = e.ID
	req.Type = pb.NotifyRequest_Type(e.Type)
	req.Code = int32(e.Code)
	req.Level = pb.NotifyRequest_Level(e.Level)
	req.Timestamp, _ = ptypes.TimestampProto(e.Timestamp)
	req.Source = &pb.NotifyRequest_Source{
		Hostname: e.Source.Hostname,
		Program:  e.Source.Program,
		Instance: e.Source.Instance,
	}
	//if no data
	if e.Data == nil || len(e.Data) == 0 {
		req.Data = &pb.NotifyRequest_Data{
			DataEnc: pb.NotifyRequest_Data_NODATA,
		}
		return req, nil
	}
	// encode data to json
	encoded, err := json.Marshal(e.Data)
	if err != nil {
		return nil, fmt.Errorf("cannot encode data to json: %v", err)
	}
	req.Data = &pb.NotifyRequest_Data{
		DataEnc: pb.NotifyRequest_Data_JSON,
		Data:    encoded,
	}
	return req, nil
}

//API returns API service name implemented
func (c *Client) API() string {
	return ServiceName()
}
