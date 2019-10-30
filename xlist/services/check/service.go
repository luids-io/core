// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package check

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/luids-io/core/protogen/xlistpb"
	"github.com/luids-io/core/xlist"
)

// Service provides a wrapper for the interface xlist.Checker that handles
// grpc requests.
type Service struct {
	opts    serviceOpts
	checker xlist.Checker
}

type serviceOpts struct {
	exposePing    bool
	disclosureErr bool
}

var defaultServiceOpts = serviceOpts{}

// ServiceOption is used for service configuration
type ServiceOption func(*serviceOpts)

// ExposePing exposes ping to the list in the service, allowing not only
// connectivity check
func ExposePing(b bool) ServiceOption {
	return func(o *serviceOpts) {
		o.exposePing = b
	}
}

// DisclosureErrors returns errors without replacing by a generic message
func DisclosureErrors(b bool) ServiceOption {
	return func(o *serviceOpts) {
		o.disclosureErr = b
	}
}

// NewService returns a new Service for the cheker
func NewService(checker xlist.Checker, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{checker: checker, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterCheckServer(server, service)
}

// Check implements grpc handler for Check
func (s *Service) Check(ctx context.Context, in *pb.CheckRequest) (*pb.CheckResponse, error) {
	name := in.GetRequest().GetName()
	resource := xlist.Resource(in.GetRequest().GetResource())
	resp, err := s.checker.Check(ctx, name, resource)
	if err != nil {
		return nil, s.mapError(err)
	}
	reply := &pb.CheckResponse{
		Response: &pb.Response{
			Result: resp.Result,
			Reason: resp.Reason,
			TTL:    int32(resp.TTL),
		}}
	return reply, nil
}

// Resources implements grpc handler for Resources
func (s *Service) Resources(ctx context.Context, in *empty.Empty) (*pb.ResourcesResponse, error) {
	resources := s.checker.Resources()
	retres := make([]pb.Resource, 0, len(resources))
	for _, r := range resources {
		retres = append(retres, pb.Resource(r))
	}
	return &pb.ResourcesResponse{Resources: retres}, nil
}

// Ping implements grpc handler for Ping
func (s *Service) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	if s.opts.exposePing {
		err := s.checker.Ping()
		if err != nil {
			return nil, s.mapError(err)
		}
	}
	return &empty.Empty{}, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	switch err {
	case xlist.ErrBadResourceFormat:
		return status.Error(codes.InvalidArgument, err.Error())
	case xlist.ErrResourceNotSupported:
		return status.Error(codes.Unimplemented, err.Error())
	case xlist.ErrListNotAvailable:
		return status.Error(codes.Unavailable, err.Error())
	default:
		rpcerr := status.Error(codes.Unavailable, xlist.ErrListNotAvailable.Error())
		if s.opts.disclosureErr {
			rpcerr = status.Error(codes.Unavailable, err.Error())
		}
		return rpcerr
	}
}
