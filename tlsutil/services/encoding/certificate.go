package encoding

import (
	"crypto/x509"

	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
)

// CertificateData convert to model CertificateData
func CertificateData(src *pb.CertificateData) *tlsutil.CertificateData {
	dst := &tlsutil.CertificateData{}
	dst.Digest = src.GetDigest()
	raw := src.GetRaw()
	if raw != nil {
		dst.Data, _ = x509.ParseCertificate(raw)
	}
	return dst
}

// CertificateDataPB convert to protobuf CertificateData
func CertificateDataPB(src *tlsutil.CertificateData) *pb.CertificateData {
	dst := &pb.CertificateData{}
	dst.Digest = src.Digest
	if src.Data != nil {
		dst.Raw = src.Data.Raw
	}
	return dst
}
