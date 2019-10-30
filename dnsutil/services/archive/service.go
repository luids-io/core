// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"net"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/dnsutil"
	pb "github.com/luids-io/core/protogen/dnsutilpb"
)

// Service implements a service wrapper for the grpc api
type Service struct {
	opts     serviceOpts
	archiver dnsutil.Archiver
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
func NewService(a dnsutil.Archiver, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{archiver: a, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterArchiveServer(server, service)
}

// SaveResolv implements interface
func (s *Service) SaveResolv(ctx context.Context, req *pb.SaveResolvRequest) (*pb.SaveResolvResponse, error) {
	//parse request
	data, err := parseRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//do request
	newid, err := s.archiver.SaveResolv(ctx, data)
	if err != nil {
		return nil, s.mapError(err)
	}
	//return response
	return &pb.SaveResolvResponse{Id: newid}, nil
}

func parseRequest(req *pb.SaveResolvRequest) (dnsutil.ResolvData, error) {
	i := dnsutil.ResolvData{}
	i.Timestamp, _ = ptypes.Timestamp(req.GetTs())
	i.Server = net.ParseIP(req.GetServerIp())
	if i.Server == nil {
		return i, dnsutil.ErrBadRequestFormat
	}
	i.Client = net.ParseIP(req.GetClientIp())
	if i.Client == nil {
		return i, dnsutil.ErrBadRequestFormat
	}
	if len(req.GetResolvedIps()) == 0 {
		return i, dnsutil.ErrBadRequestFormat
	}
	i.Resolved = make([]net.IP, 0, len(req.GetResolvedIps()))
	for _, r := range req.GetResolvedIps() {
		ip := net.ParseIP(r)
		if ip == nil {
			return i, dnsutil.ErrBadRequestFormat
		}
		i.Resolved = append(i.Resolved, ip)
	}
	i.Name = req.GetName()
	if i.Name == "" {
		return i, dnsutil.ErrBadRequestFormat
	}
	return i, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	switch err {
	case dnsutil.ErrBadRequestFormat:
		return status.Error(codes.InvalidArgument, err.Error())
	case dnsutil.ErrServiceNotAvailable:
		return status.Error(codes.Unavailable, err.Error())
	default:
		rpcerr := status.Error(codes.Unavailable, dnsutil.ErrServiceNotAvailable.Error())
		if s.opts.disclosureErr {
			rpcerr = status.Error(codes.Unavailable, err.Error())
		}
		return rpcerr
	}
}
