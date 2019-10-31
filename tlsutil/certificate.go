// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"crypto/x509"
)

// CertificateData stores certificate information
type CertificateData struct {
	ID     string            `json:"id"`
	Digest string            `json:"digest"`
	Data   *x509.Certificate `json:"data"`
}

// CertSummary stores basic information of certification
type CertSummary struct {
	Digest  string `json:"digest"`
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
	IsCA    bool   `json:"isCA"`
}

// Summary returns certificate summary
func (c *CertificateData) Summary() CertSummary {
	return CertSummary{
		Digest:  c.Digest,
		Issuer:  c.Data.Issuer.CommonName,
		Subject: c.Data.Subject.CommonName,
		IsCA:    c.Data.IsCA,
	}
}
