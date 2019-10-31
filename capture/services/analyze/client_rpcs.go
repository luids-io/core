// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package analyze

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/gopacket/layers"
	"google.golang.org/grpc"

	pb "github.com/luids-io/core/protogen/capturepb"
)

type pcktClientStream interface {
	grpc.ClientStream
	Send(*pb.SendPacketRequest) error
}

type rpcClient struct {
	client    pb.AnalyzeClient
	stream    pcktClientStream
	dataCh    chan *pb.SendPacketRequest
	linkType  layers.LinkType
	connected bool
}

func newRPCClient(c pb.AnalyzeClient, l layers.LinkType, buffSize int) *rpcClient {
	r := &rpcClient{
		client:   c,
		linkType: l,
		dataCh:   make(chan *pb.SendPacketRequest, buffSize),
	}
	return r
}

// Data returns channel for write data
func (r *rpcClient) Data() chan<- *pb.SendPacketRequest {
	return r.dataCh
}

func (r *rpcClient) run(wg *sync.WaitGroup, closeCh <-chan struct{}, errCh chan<- error) {
PROCESSLOOP:
	for {
		select {
		case data := <-r.dataCh:
			err := r.save(data)
			if err != nil {
				errCh <- err
			}
		case <-closeCh:
			//clean buffer
			for data := range r.dataCh {
				err := r.save(data)
				if err != nil {
					errCh <- err
				}
			}
			break PROCESSLOOP
		}
	}
	//close channel data and close stream
	close(r.dataCh)
	r.close()

	wg.Done()
}

//save request, implements a reconnection system
func (r *rpcClient) save(req *pb.SendPacketRequest) error {
	if !r.connected {
		err := r.connect()
		if err != nil {
			return err
		}
	}
	//send
	err := r.stream.Send(req)
	if err != nil {
		r.close()
	}
	return err
}

func (r *rpcClient) connect() error {
	var err error
	switch r.linkType {
	case layers.LinkTypeEthernet:
		r.stream, err = r.client.SendEtherPackets(context.Background())
	default:
		err = fmt.Errorf("unexpected linktype %v", r.linkType)
	}
	if err != nil {
		return err
	}
	r.connected = true
	return nil
}

func (r *rpcClient) close() {
	if r.connected {
		r.stream.CloseSend()
		r.connected = false
	}
}
