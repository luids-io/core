// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/event"
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
	e.Received = time.Now()
	reqID, err := s.notifier.Notify(ctx, e)
	if err != nil {
		rpcerr := status.Error(codes.Internal, err.Error())
		return nil, rpcerr
	}
	reply := &pb.NotifyResponse{RequestID: reqID}
	return reply, nil
}

func eventFromRequest(req *pb.NotifyRequest) (event.Event, error) {
	e := event.Event{}
	e.ID = req.GetId()
	e.Type = event.Type(req.GetType())
	e.Code = event.Code(req.GetCode())
	e.Level = event.Level(req.GetLevel())
	e.Timestamp, _ = ptypes.Timestamp(req.GetTimestamp())
	e.Source.Hostname = req.GetSource().GetHostname()
	e.Source.Instance = req.GetSource().GetInstance()
	e.Source.Program = req.GetSource().GetProgram()
	//decode event data
	switch req.GetData().GetDataEnc() {
	case pb.NotifyRequest_Data_JSON:
		rawData := req.Data.GetData()
		err := json.Unmarshal(rawData, &e.Data)
		if err != nil {
			return event.Event{}, fmt.Errorf("unmarshalling data: %v", err)
		}
	case pb.NotifyRequest_Data_NODATA:
		e.Data = make(map[string]interface{})
	}
	return e, nil
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}
