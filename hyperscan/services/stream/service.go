// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package stream

import (
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/hyperscan"
	pb "github.com/luids-io/core/protogen/hspb"
)

// Service provides a wrapper for the interface hyperscan.BlockScanner that handles
// grpc requests.
type Service struct {
	opts    serviceOpts
	scanner hyperscan.StreamScanner
}

type serviceOpts struct {
	disclosureErr bool
	dataBuff      int
}

var defaultServiceOpts = serviceOpts{dataBuff: 100}

// ServiceOption is used for service configuration
type ServiceOption func(*serviceOpts)

// DisclosureErrors returns errors without replacing by a generic message
func DisclosureErrors(b bool) ServiceOption {
	return func(o *serviceOpts) {
		o.disclosureErr = b
	}
}

// SetDataBuff option allows change channel buffer data
func SetDataBuff(i int) ServiceOption {
	return func(o *serviceOpts) {
		o.dataBuff = i
	}
}

// NewService returns a new Service for the cheker
func NewService(s hyperscan.StreamScanner, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{scanner: s, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterStreamServer(server, service)
}

// ScanStream implements interface
func (s *Service) ScanStream(stream pb.Stream_ScanStreamServer) error {
	ctx := stream.Context()

	dataCh := make(chan []byte, s.opts.dataBuff)
	responses, err := s.scanner.ScanStream(ctx, dataCh)
	if err != nil {
		return status.Errorf(codes.Internal, "Internal error")
	}

	closed := make(chan struct{})
	// send responses go routine
	go func() {
	SENDLOOP:
		for {
			select {
			case <-ctx.Done():
				break SENDLOOP
			case reason, ok := <-responses:
				if !ok {
					break SENDLOOP
				}
				stream.Send(&pb.ScanStreamResponse{Reason: reason})
			}
		}
		close(closed)
	}()
	// receive data go routine
	go func() {
	RCVDLOOP:
		for {
			request, err := stream.Recv()
			if err == io.EOF {
				break RCVDLOOP
			} else if err != nil {
				break RCVDLOOP
			}
			if request.GetFinish() {
				break RCVDLOOP
			}
			data := request.GetData()
			if data == nil {
				continue
			}
			dataCh <- data
		}
		close(dataCh)
	}()

	//wait for close
	<-closed
	return nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	//TODO
	rpcerr := status.Error(codes.Unavailable, "service not available")
	if s.opts.disclosureErr {
		rpcerr = status.Error(codes.Unavailable, err.Error())
	}
	return rpcerr
}
