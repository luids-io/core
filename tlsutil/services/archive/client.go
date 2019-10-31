// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
)

// Client is the main struct for grpc client
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.ArchiveClient
	//stream rpc management
	recordRPC *rpcRecord
	//control
	started bool
	close   chan struct{}
	errs    chan error
	wg      sync.WaitGroup
}

type clientOpts struct {
	logger    yalogi.Logger
	closeConn bool
	buffSize  int
}

var defaultClientOpts = clientOpts{
	logger:    yalogi.LogNull,
	closeConn: true,
	buffSize:  100,
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
	c := &Client{
		opts:   opts,
		logger: opts.logger,
		conn:   conn,
		client: pb.NewArchiveClient(conn),
	}
	c.start()
	return c
}

// SaveConnection implements tlsutil.Archiver interface
func (c *Client) SaveConnection(ctx context.Context, data *tlsutil.ConnectionData) (string, error) {
	if !c.started {
		return "", errors.New("client closed")
	}
	req, err := connectionToRequest(data)
	if err != nil {
		return "", c.mapError(err)
	}
	resp, err := c.client.SaveConnection(ctx, req)
	if err != nil {
		return "", c.mapError(err)
	}
	return resp.GetId(), nil
}

// SaveCertificate implements tlsutil.Archiver interface
func (c *Client) SaveCertificate(ctx context.Context, data *tlsutil.CertificateData) (string, error) {
	if !c.started {
		return "", errors.New("client closed")
	}
	req, err := certificateToRequest(data)
	if err != nil {
		return "", c.mapError(err)
	}
	resp, err := c.client.SaveCertificate(ctx, req)
	if err != nil {
		return "", c.mapError(err)
	}
	return resp.GetId(), nil
}

// StoreRecord implements tlsutil.Archiver interface
func (c *Client) StoreRecord(data *tlsutil.RecordData) error {
	if !c.started {
		return errors.New("client closed")
	}
	req := recordToRequest(data)
	//TODO: check errs
	c.recordRPC.Data() <- req
	return nil
}

func (c *Client) start() {
	//init status
	c.close = make(chan struct{})
	c.errs = make(chan error, c.opts.buffSize)
	go c.processErrs()

	//init rpc managers
	c.wg.Add(1)
	c.recordRPC = newRPCrecord(c.client, c.opts.buffSize)
	go c.recordRPC.run(&c.wg, c.close, c.errs)

	c.started = true
}

func (c *Client) processErrs() {
	for e := range c.errs {
		c.logger.Warnf("%v", e)
	}
}

//Close closes the client
func (c *Client) Close() error {
	if !c.started {
		return errors.New("client closed")
	}
	c.started = false
	//close signal
	close(c.close)
	c.wg.Wait()
	close(c.errs)

	if c.opts.closeConn {
		return c.conn.Close()
	}
	return nil
}

//mapping errors routine
func (c *Client) mapError(err error) error {
	//TODO
	return err
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
