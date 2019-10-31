// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package analyze

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/luids-io/core/capture"
	pb "github.com/luids-io/core/protogen/capturepb"
)

// Service implements a service wrapper for the grpc api
type Service struct {
	processor capture.Processor
}

// NewService returns a new Service for the grpc api
func NewService(p capture.Processor) *Service {
	return &Service{processor: p}
}

// RegisterServer registers a service in the grpc server
func RegisterServer(server *grpc.Server, service *Service) {
	pb.RegisterAnalyzeServer(server, service)
}

// SendEtherPackets manage requests
func (s *Service) SendEtherPackets(stream pb.Analyze_SendEtherPacketsServer) error {
	return s.sendPackets(stream, layers.LinkTypeEthernet)
}

// sendPackets manage requests
func (s *Service) sendPackets(stream pcktServerStream, linkType layers.LinkType) error {
	ctx := stream.Context()
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Errorf(codes.Internal, "Internal error getting peer")
	}
	// creates packet source
	name := p.Addr.String()
	switch linkType {
	case layers.LinkTypeEthernet:
		name = fmt.Sprintf("%s-eth", name)
	default:
		return status.Error(codes.Internal, "invalid layer")
	}
	psource := &pcktSource{
		err:    make(chan error),
		stream: stream,
	}
	source := gopacket.NewPacketSource(psource, linkType)
	err := s.processor.Register(name, source)
	if err != nil {
		return status.Errorf(codes.Internal, "Internal error registering")
	}
	//waits for close or error
	err = <-psource.err
	//clean
	s.processor.Unregister(name)
	close(psource.err)
	if err == io.EOF {
		return nil
	}
	return err
}

//mapping errors
func (s *Service) mapError(err error) error {
	return status.Error(codes.Unavailable, err.Error())
}

type pcktSource struct {
	stream pcktServerStream
	err    chan error
	closed bool
}

type pcktServerStream interface {
	grpc.ServerStream
	Recv() (*pb.SendPacketRequest, error)
}

func (p *pcktSource) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if p.closed {
		err = io.EOF
		return
	}
	var req *pb.SendPacketRequest
	req, err = p.stream.Recv()
	if err != nil {
		p.closed = true
		p.err <- err
		return
	}
	meta := req.GetMetadata()
	ts := meta.GetTimestamp()
	data = req.GetData()
	ci.Timestamp, _ = ptypes.Timestamp(ts)
	ci.InterfaceIndex = int(meta.GetInterfaceIndex())
	ci.CaptureLength = len(data)
	ci.Length = len(data)
	return
}
