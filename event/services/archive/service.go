// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

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
	archiver event.Archiver
}

// NewService returns a new Service for the grpc api
func NewService(a event.Archiver) *Service {
	return &Service{archiver: a}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterArchiveServer(server, service)
}

// SaveEvent implements API service
func (s *Service) SaveEvent(ctx context.Context, in *pb.SaveEventRequest) (*pb.SaveEventResponse, error) {
	e, err := eventFromRequest(in)
	if err != nil {
		rpcerr := status.Error(codes.InvalidArgument, "request is not valid")
		return nil, rpcerr
	}
	eID, err := s.archiver.SaveEvent(ctx, e)
	if err != nil {
		rpcerr := status.Error(codes.Internal, err.Error())
		return nil, rpcerr
	}
	reply := &pb.SaveEventResponse{EventID: eID}
	return reply, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}

func eventFromRequest(req *pb.SaveEventRequest) (event.Event, error) {
	pbevent := req.GetEvent()
	if pbevent == nil {
		return event.Event{}, errors.New("event is empty")
	}
	return encoding.Event(pbevent)
}
