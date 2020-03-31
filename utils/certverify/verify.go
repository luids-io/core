package certverify

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"
)

var rootsCA *x509.CertPool

func caBundle(caBundlePath string) (*x509.CertPool, error) {
	if caBundlePath == "" {
		return x509.SystemCertPool()
	}

	caBundleBytes, err := ioutil.ReadFile(caBundlePath)
	if err != nil {
		return nil, err
	}

	bundle := x509.NewCertPool()
	ok := bundle.AppendCertsFromPEM(caBundleBytes)
	if !ok {
		return nil, errors.New("unable to read certificates from CA bundle")
	}

	return bundle, nil
}

// SetCABundlePath sets file path where CA certs are stored
func SetCABundlePath(caBundlePath string) error {
	var err error
	rootsCA, err = caBundle(caBundlePath)

	return err
}

// VerifyChain verifies a chain of certification passed in a slice
func VerifyChain(certs []*x509.Certificate, currentTime time.Time, dnsName string) error {

	intermediates := x509.NewCertPool()
	for i := 1; i < len(certs); i++ {
		intermediates.AddCert(certs[i])
	}

	opts := x509.VerifyOptions{
		DNSName:       dnsName,
		Roots:         rootsCA,
		CurrentTime:   currentTime,
		Intermediates: intermediates,
	}

	_, err := certs[0].Verify(opts)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootsCA, _ = caBundle("")
}
