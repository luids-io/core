// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package block

import (
	"context"

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
	scanner hyperscan.BlockScanner
}

type serviceOpts struct {
	disclosureErr bool
}

var defaultServiceOpts = serviceOpts{}

// ServiceOption is used for service configuration
type ServiceOption func(*serviceOpts)

// DisclosureErrors returns errors without replacing by a generic message
func DisclosureErrors(b bool) ServiceOption {
	return func(o *serviceOpts) {
		o.disclosureErr = b
	}
}

// NewService returns a new Service for the cheker
func NewService(s hyperscan.BlockScanner, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{scanner: s, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterBlockServer(server, service)
}

// ScanBlock implements grpc handler for Scan
func (s *Service) ScanBlock(ctx context.Context, in *pb.ScanBlockRequest) (*pb.ScanBlockResponse, error) {
	data := in.GetData()
	matches, reasons, err := s.scanner.ScanBlock(ctx, data)
	if err != nil {
		return nil, s.mapError(err)
	}
	return &pb.ScanBlockResponse{Result: matches, Reasons: reasons}, nil
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
