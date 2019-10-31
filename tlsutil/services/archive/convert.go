// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package archive

import (
	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
	"github.com/luids-io/core/tlsutil/services/encoding"
)

func certificateToRequest(cert *tlsutil.CertificateData) (*pb.SaveCertificateRequest, error) {
	req := &pb.SaveCertificateRequest{}
	req.Certificate = encoding.CertificateDataPB(cert)
	return req, nil
}

func certificateFromRequest(req *pb.SaveCertificateRequest) (*tlsutil.CertificateData, error) {
	var cert *tlsutil.CertificateData
	c := req.GetCertificate()
	if c != nil {
		cert = encoding.CertificateData(c)
	}
	return cert, nil
}

func connectionToRequest(cn *tlsutil.ConnectionData) (*pb.SaveConnectionRequest, error) {
	req := &pb.SaveConnectionRequest{}
	req.Connection = encoding.ConnectionDataPB(cn)
	return req, nil
}

func connectionFromRequest(req *pb.SaveConnectionRequest) (*tlsutil.ConnectionData, error) {
	var cn *tlsutil.ConnectionData
	c := req.GetConnection()
	if c != nil {
		cn = encoding.ConnectionData(c)
	}
	return cn, nil
}

func recordToRequest(r *tlsutil.RecordData) *pb.SaveRecordRequest {
	req := &pb.SaveRecordRequest{}
	req.Record = encoding.RecordDataPB(r)
	return req
}

func recordFromRequest(req *pb.SaveRecordRequest) *tlsutil.RecordData {
	var record *tlsutil.RecordData
	r := req.GetRecord()
	if r != nil {
		record = encoding.RecordData(r)
	}
	return record
}
