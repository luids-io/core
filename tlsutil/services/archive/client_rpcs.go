// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	"context"
	"sync"

	"github.com/luids-io/core/protogen/tlsutilpb"
	pb "github.com/luids-io/core/protogen/tlsutilpb"
)

// rpcRecord manage stream for sending records
type rpcRecord struct {
	client    pb.ArchiveClient
	stream    pb.Archive_StreamRecordsClient
	dataCh    chan *pb.SaveRecordRequest
	connected bool
}

func newRPCrecord(c pb.ArchiveClient, buffSize int) *rpcRecord {
	r := &rpcRecord{
		client: c,
		dataCh: make(chan *pb.SaveRecordRequest, buffSize),
	}
	return r
}

// Data returns channel for write data
func (r *rpcRecord) Data() chan<- *pb.SaveRecordRequest {
	return r.dataCh
}

func (r *rpcRecord) run(wg *sync.WaitGroup, closeCh <-chan struct{}, errCh chan<- error) {
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

func (r *rpcRecord) save(req *tlsutilpb.SaveRecordRequest) error {
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

func (r *rpcRecord) connect() error {
	var err error
	r.stream, err = r.client.StreamRecords(context.Background())
	if err != nil {
		return err
	}
	r.connected = true
	return nil
}

func (r *rpcRecord) close() {
	if r.connected {
		r.stream.CloseAndRecv()
		r.connected = false
	}
}
