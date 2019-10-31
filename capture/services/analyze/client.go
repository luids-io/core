// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package analyze

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	pb "github.com/luids-io/core/protogen/capturepb"
)

// Client is the main struct for grpc client
type Client struct {
	opts   clientOpts
	logger yalogi.Logger
	//grpc connection
	conn   *grpc.ClientConn
	client pb.AnalyzeClient
	//rpc management
	rpcEth *rpcClient
	rpcIP4 *rpcClient
	rpcIP6 *rpcClient
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
		client: pb.NewAnalyzeClient(conn),
	}
	c.start()
	return c
}

// SendEtherPacket implements capture.Analyzer interface
func (c *Client) SendEtherPacket(packet gopacket.Packet) error {
	if !c.started {
		return errors.New("client closed")
	}
	req, err := getPacketRequest(packet, layers.LayerTypeEthernet)
	if err != nil {
		return err
	}
	c.rpcEth.Data() <- req
	return nil
}

// SendIP4Packet implements capture.Analyzer interface
func (c *Client) SendIP4Packet(packet gopacket.Packet) error {
	if !c.started {
		return errors.New("client closed")
	}
	req, err := getPacketRequest(packet, layers.LayerTypeIPv4)
	if err != nil {
		return err
	}
	c.rpcIP4.Data() <- req
	return nil
}

// SendIP6Packet implements capture.Analyzer interface
func (c *Client) SendIP6Packet(packet gopacket.Packet) error {
	if !c.started {
		return errors.New("client closed")
	}
	req, err := getPacketRequest(packet, layers.LayerTypeIPv6)
	if err != nil {
		return err
	}
	c.rpcIP6.Data() <- req
	return nil
}

func getPacketRequest(packet gopacket.Packet, layer gopacket.LayerType) (*pb.SendPacketRequest, error) {
	ldata := packet.Layer(layer)
	if ldata == nil {
		return nil, errors.New("layer missmatch")
	}
	//create request
	md := packet.Metadata()
	req := &pb.SendPacketRequest{}
	req.Metadata = &pb.PacketMetadata{}
	req.Metadata.Timestamp, _ = ptypes.TimestampProto(md.Timestamp)
	req.Metadata.InterfaceIndex = int32(md.InterfaceIndex)
	req.Data = append(ldata.LayerContents(), ldata.LayerPayload()...)
	return req, nil
}

func (c *Client) start() {
	//init status
	c.close = make(chan struct{})
	c.errs = make(chan error, c.opts.buffSize)
	go c.processErrs()

	//init rpc managers
	c.wg.Add(1)
	c.rpcEth = newRPCClient(c.client, layers.LinkTypeEthernet, c.opts.buffSize)
	go c.rpcEth.run(&c.wg, c.close, c.errs)

	c.wg.Add(1)
	c.rpcIP4 = newRPCClient(c.client, layers.LinkTypeIPv4, c.opts.buffSize)
	go c.rpcIP4.run(&c.wg, c.close, c.errs)

	c.wg.Add(1)
	c.rpcIP6 = newRPCClient(c.client, layers.LinkTypeIPv6, c.opts.buffSize)
	go c.rpcIP6.run(&c.wg, c.close, c.errs)

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
