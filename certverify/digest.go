package certverify

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
)

// DigestCert returns a digest of a certificate
func DigestCert(cert *x509.Certificate) string {
	hasher := md5.New()
	return hex.EncodeToString(hasher.Sum(cert.Raw))
}

// DigestChain returns digest of a certificate chain
func DigestChain(certs []*x509.Certificate) string {
	digest := ""
	for _, cert := range certs {
		digest += DigestCert(cert)
	}
	return hashString(digest)
}

func hashString(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
