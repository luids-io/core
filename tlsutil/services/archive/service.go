// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
)

// Service implements a service wrapper for the grpc api
type Service struct {
	archiver tlsutil.Archiver
}

// NewService returns a new Service for the grpc api
func NewService(a tlsutil.Archiver) *Service {
	return &Service{archiver: a}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterArchiveServer(server, service)
}

// SaveConnection implements interface
func (s *Service) SaveConnection(ctx context.Context, req *pb.SaveConnectionRequest) (*pb.SaveConnectionResponse, error) {
	//parse request
	data, err := connectionFromRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//do request
	newid, err := s.archiver.SaveConnection(ctx, data)
	if err != nil {
		return nil, s.mapError(err)
	}
	//return response
	return &pb.SaveConnectionResponse{Id: newid}, nil
}

// SaveCertificate implements interface
func (s *Service) SaveCertificate(ctx context.Context, req *pb.SaveCertificateRequest) (*pb.SaveCertificateResponse, error) {
	//parse request
	data, err := certificateFromRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	//do request
	newid, err := s.archiver.SaveCertificate(ctx, data)
	if err != nil {
		return nil, s.mapError(err)
	}
	//return response
	return &pb.SaveCertificateResponse{Id: newid}, nil
}

// StreamRecords implements interface
func (s *Service) StreamRecords(stream pb.Archive_StreamRecordsServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		record := recordFromRequest(request)
		err = s.archiver.StoreRecord(record)
		if err != nil {
			return err
		}
	}
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}
