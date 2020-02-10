// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package classify

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/brain/classify"
	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil/services/encoding"
)

// Service provides a grpc wrapper
type Service struct {
	opts       serviceOpts
	classifier classify.Classifier
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
func NewService(c classify.Classifier, opt ...ServiceOption) *Service {
	opts := defaultServiceOpts
	for _, o := range opt {
		o(&opts)
	}
	return &Service{classifier: c, opts: opts}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterClassifyServer(server, service)
}

// Connections implements grpc interface
func (s *Service) Connections(ctx context.Context, in *pb.ClassifyConnectionRequest) (*pb.ClassifyConnectionResponse, error) {
	// prepare request
	if len(in.Requests) == 0 {
		rpcerr := status.Error(codes.InvalidArgument, "requests is empty")
		return nil, rpcerr
	}
	requests := make([]classify.Request, 0, len(in.Requests))
	for _, r := range in.Requests {
		requests = append(requests, classify.Request{
			ID:   r.GetId(),
			Data: encoding.ConnectionData(r.GetConnection()),
		})
	}
	// do request
	responses, err := s.classifier.Classify(ctx, requests)
	if err != nil {
		return nil, s.mapError(err)
	}
	// create response
	retResponses := make([]*pb.ClassifyConnectionResponse_Response, 0, len(responses))
	for _, r := range responses {
		resp := &pb.ClassifyConnectionResponse_Response{Id: r.ID}
		retResponses = append(retResponses, resp)
		if r.Err != nil {
			resp.Err = r.Err.Error()
			continue
		}
		resp.Results = make([]*pb.ClassifyConnectionResponse_Response_Result, 0, len(r.Results))
		for _, result := range r.Results {
			resp.Results = append(resp.Results, &pb.ClassifyConnectionResponse_Response_Result{
				Label: result.Label,
				Prob:  result.Prob},
			)
		}
	}
	return &pb.ClassifyConnectionResponse{Responses: retResponses}, nil
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
