// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"errors"
	"fmt"

	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/luids-io/core/event"
	"github.com/luids-io/core/event/services/encoding"
	pb "github.com/luids-io/core/protogen/eventpb"
)

// Client is the main struct for grpc client
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

// SaveEvent implements event.Archiver interface
func (c *Client) SaveEvent(ctx context.Context, e event.Event) (string, error) {
	if !c.started {
		return "", errors.New("client closed")
	}
	return c.doSaveEvent(ctx, e)
}

func (c *Client) doSaveEvent(ctx context.Context, e event.Event) (string, error) {
	//create request
	req, err := eventToRequest(e)
	if err != nil {
		return "", fmt.Errorf("serializing event: %v", err)
	}
	//notify request
	resp, err := c.client.SaveEvent(ctx, req)
	if err != nil {
		return "", c.mapError(err)
	}
	//process response
	eventID := resp.GetEventID()
	if eventID == "" {
		return "", fmt.Errorf("processing response: event_id empty")
	}
	return eventID, nil
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

//API returns API service name implemented
func (c *Client) API() string {
	return ServiceName()
}

func eventToRequest(e event.Event) (*pb.SaveEventRequest, error) {
	pbevent, err := encoding.EventPB(e)
	if err != nil {
		return nil, err
	}
	req := &pb.SaveEventRequest{}
	req.Event = pbevent
	return req, nil
}
