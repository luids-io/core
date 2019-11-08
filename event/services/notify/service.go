// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package notify

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/event"
	"github.com/luids-io/core/event/services/encoding"
	pb "github.com/luids-io/core/protogen/eventpb"
)

// Service implements a service wrapper for the grpc api
type Service struct {
	notifier event.Notifier
}

// NewService returns a new Service for the grpc api
func NewService(n event.Notifier) *Service {
	return &Service{notifier: n}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterNotifyServer(server, service)
}

// Notify implements API service
func (s *Service) Notify(ctx context.Context, in *pb.NotifyRequest) (*pb.NotifyResponse, error) {
	e, err := eventFromRequest(in)
	if err != nil {
		rpcerr := status.Error(codes.InvalidArgument, "request is not valid")
		return nil, rpcerr
	}
	reqID, err := s.notifier.Notify(ctx, e)
	if err != nil {
		rpcerr := status.Error(codes.Internal, err.Error())
		return nil, rpcerr
	}
	reply := &pb.NotifyResponse{RequestID: reqID}
	return reply, nil
}

func eventFromRequest(req *pb.NotifyRequest) (event.Event, error) {
	pbevent := req.GetEvent()
	if pbevent == nil {
		return event.Event{}, errors.New("event is empty")
	}
	return encoding.Event(pbevent)
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}
