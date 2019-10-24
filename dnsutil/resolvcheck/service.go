// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package resolvcheck

import (
	"context"
	"errors"
	"net"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/dnsutil"
	pb "github.com/luids-io/core/protogen/dnsutilpb"
)

// Service provides a wrapper
type Service struct {
	checker dnsutil.ResolvChecker
}

// NewService returns a new Service
func NewService(c dnsutil.ResolvChecker) *Service {
	return &Service{checker: c}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterResolvCheckServer(server, service)
}

// Check implements API
func (s *Service) Check(ctx context.Context, in *pb.ResolvCheckRequest) (*pb.ResolvCheckResponse, error) {
	//parse request
	client, resolved, name, err := parseRequest(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//do request
	resp, err := s.checker.Check(ctx, client, resolved, name)
	if err != nil {
		return nil, s.mapError(err)
	}
	//return response
	response := &pb.ResolvCheckResponse{}
	response.Result = resp.Result
	response.LastTs, _ = ptypes.TimestampProto(resp.Last)
	response.StoreTs, _ = ptypes.TimestampProto(resp.Store)
	return response, nil
}

func parseRequest(req *pb.ResolvCheckRequest) (net.IP, net.IP, string, error) {
	client := req.GetClientIp()
	resolved := req.GetResolvedIp()
	if client == "" || resolved == "" {
		return nil, nil, "", errors.New("client and resolved are required")
	}
	clientIP := net.ParseIP(client)
	if clientIP == nil {
		return nil, nil, "", errors.New("client must be an ip")
	}
	resolvedIP := net.ParseIP(resolved)
	if resolvedIP == nil {
		return nil, nil, "", errors.New("resolved must be an ip")
	}
	name := req.GetName()
	return clientIP, resolvedIP, name, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}
