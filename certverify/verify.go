// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package certverify provides some helper functions to verify, compute the
// signature or download certificates.
//
// This package is a work in progress and makes no API stability promises.
package certverify

import (
	"crypto/x509"
	"errors"
	"io/ioutil"
	"time"
)

// VerifyChain verifies a chain of certificates passed in a slice.
func VerifyChain(certs []*x509.Certificate, currentTime time.Time, dnsName string) error {
	// prepare itermediates
	intermediates := x509.NewCertPool()
	for i := 1; i < len(certs); i++ {
		intermediates.AddCert(certs[i])
	}
	// define options
	opts := x509.VerifyOptions{
		DNSName:       dnsName,
		Roots:         rootsCA,
		CurrentTime:   currentTime,
		Intermediates: intermediates,
	}
	//verify certs
	_, err := certs[0].Verify(opts)
	if err != nil {
		return err
	}
	return nil
}

// SetCABundlePath sets file path where CA certs are stored.
func SetCABundlePath(caBundlePath string) error {
	rca, err := caBundle(caBundlePath)
	if err != nil {
		return err
	}
	rootsCA = rca
	return nil
}

var rootsCA *x509.CertPool

func caBundle(caBundlePath string) (*x509.CertPool, error) {
	caBundleBytes, err := ioutil.ReadFile(caBundlePath)
	if err != nil {
		return nil, err
	}
	//create pool
	bundle := x509.NewCertPool()
	ok := bundle.AppendCertsFromPEM(caBundleBytes)
	if !ok {
		return nil, errors.New("unable to read certificates from CA bundle")
	}
	return bundle, nil
}
