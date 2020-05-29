// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package certverify

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

// FetchCerts from a tls server doing a nil connection, returning certificates and oscp response
func FetchCerts(ip net.IP, port int, sni string, config *tls.Config) ([]*x509.Certificate, []byte, error) {
	config = config.Clone()
	// prepare config to fetch certs
	config.ServerName = sni
	config.ClientSessionCache = &nilSessionCache{}
	config.InsecureSkipVerify = true
	// do connection
	conn, err := tls.Dial("tcp", fmt.Sprintf("%v:%v", ip, port), config)
	if err != nil {
		return nil, nil, err
	}
	// get certs and oscp
	certs := conn.ConnectionState().PeerCertificates
	oscp := conn.ConnectionState().OCSPResponse
	return certs, oscp, nil
}

// implemented a nil session because we don't want to reuse connections
type nilSessionCache struct{}

func (c *nilSessionCache) Get(sessionKey string) (session *tls.ClientSessionState, ok bool) {
	return nil, false
}
func (c *nilSessionCache) Put(sessionKey string, cs *tls.ClientSessionState) {}
