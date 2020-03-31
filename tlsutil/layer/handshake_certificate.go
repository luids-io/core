// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"crypto/x509"
	"fmt"
)

// CertificateData is the struct for protocol hanshake message Certificate
type CertificateData struct {
	CertificatesLen uint32              `json:"certificatesLen"`
	Certificates    []*x509.Certificate `json:"certificates,omitempty"`
}

func (hs *CertificateData) String() string {
	str := fmt.Sprintln("Certificates Len:", hs.CertificatesLen)

	return str
}
