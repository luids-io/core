// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package resolvcollect

import (
	"context"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/dnsutil"
	pb "github.com/luids-io/core/protogen/dnsutilpb"
)

// Service implements a service wrapper for the grpc api
type Service struct {
	opts      serviceOpts
	collector dnsutil.ResolvCollector
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

// NewService returns a new Service for the grpc api
func NewService(c dnsutil.ResolvCollector, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{collector: c, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterResolvCollectServer(server, service)
}

// Collect implements interface
func (s *Service) Collect(ctx context.Context, req *pb.ResolvCollectRequest) (*empty.Empty, error) {
	//parse request
	client, name, resolved, err := parseRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//do request
	err = s.collector.Collect(ctx, client, name, resolved)
	if err != nil {
		return nil, s.mapError(err)
	}
	//return response
	return &empty.Empty{}, nil
}

func parseRequest(req *pb.ResolvCollectRequest) (net.IP, string, []net.IP, error) {
	client := net.ParseIP(req.GetClientIp())
	if client == nil {
		return nil, "", nil, dnsutil.ErrBadRequestFormat
	}
	name := req.GetName()
	if name == "" {
		return nil, "", nil, dnsutil.ErrBadRequestFormat
	}
	if len(req.GetResolvedIps()) == 0 {
		return nil, "", nil, dnsutil.ErrBadRequestFormat
	}
	resolved := make([]net.IP, 0, len(req.GetResolvedIps()))
	for _, r := range req.GetResolvedIps() {
		ip := net.ParseIP(r)
		if ip == nil {
			return nil, "", nil, dnsutil.ErrBadRequestFormat
		}
		resolved = append(resolved, ip)
	}
	return client, name, resolved, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	switch err {
	case dnsutil.ErrBadRequestFormat:
		return status.Error(codes.InvalidArgument, err.Error())
	case dnsutil.ErrCollectDNSClientLimit:
		return status.Error(codes.ResourceExhausted, err.Error())
	case dnsutil.ErrCollectNamesLimit:
		return status.Error(codes.ResourceExhausted, err.Error())
	case dnsutil.ErrCacheNotAvailable:
		return status.Error(codes.Unavailable, err.Error())
	default:
		rpcerr := status.Error(codes.Unavailable, dnsutil.ErrCacheNotAvailable.Error())
		if s.opts.disclosureErr {
			rpcerr = status.Error(codes.Unavailable, err.Error())
		}
		return rpcerr
	}
}
