// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package certverify

import (
	"crypto/x509"
	"encoding/hex"
	"hash"
)

// DigestCert returns a hash of the certificate encoded in string
func DigestCert(hasher hash.Hash, cert *x509.Certificate) string {
	hasher.Reset()
	hasher.Write(cert.Raw)
	return hex.EncodeToString(hasher.Sum(nil))
}

// DigestCerts returns an slice with hashes of certificates and a hash of the concatenation
func DigestCerts(hasher hash.Hash, certs []*x509.Certificate) ([]string, string) {
	chain := ""
	hashes := make([]string, 0, len(certs))
	for _, cert := range certs {
		digest := DigestCert(hasher, cert)
		hashes = append(hashes, digest)
		chain += digest
	}
	hasher.Reset()
	hasher.Write([]byte(chain))
	return hashes, hex.EncodeToString(hasher.Sum(nil))
}

// DigestChain returns a hash of the concatenation of digests
func DigestChain(hasher hash.Hash, digests []string) string {
	chain := ""
	for _, digest := range digests {
		chain += digest
	}
	hasher.Reset()
	hasher.Write([]byte(chain))
	return hex.EncodeToString(hasher.Sum(nil))
}
