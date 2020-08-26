// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package grpctls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"google.golang.org/grpc/credentials"
)

// ServerCfg defines configuration for a server.
type ServerCfg struct {
	// CertFile path to the server certificate
	CertFile string
	// KeyFile path to the private key
	KeyFile string
	// CACert path to certification authority certificate
	CACert string
	// ClientAuth if client certificate is required
	ClientAuth bool
}

// UseTLS returns true if TLS configuration is set.
func (cfg ServerCfg) UseTLS() bool {
	if cfg.CertFile == "" {
		return false
	}
	return true
}

// Validate if configuration is ok.
func (cfg ServerCfg) Validate() error {
	if cfg.CertFile == "" {
		return errors.New("certfile is required")
	}
	if !fileExists(cfg.CertFile) {
		return fmt.Errorf("certfile '%v' doesn't exists", cfg.CertFile)
	}
	if cfg.KeyFile == "" {
		return errors.New("keyfile is required for server's certificate")
	}
	if !fileExists(cfg.KeyFile) {
		return fmt.Errorf("keyfile '%v' doesn't exists", cfg.KeyFile)
	}
	if cfg.CACert != "" {
		if !fileExists(cfg.CACert) {
			return fmt.Errorf("cacert file '%v' doesn't exists", cfg.CACert)
		}
	}
	if cfg.ClientAuth {
		switch {
		case cfg.CertFile == "":
			return errors.New("cert file is required for client auth")
		case cfg.KeyFile == "":
			return errors.New("key file is required for client auth")
		case cfg.CACert == "":
			return errors.New("cacert file is required for client auth")
		}
	}
	return nil
}

// Creds creates a transport credentials for the configuration.
func Creds(cfg ServerCfg) (credentials.TransportCredentials, error) {
	// some checks
	if !cfg.UseTLS() {
		return nil, errors.New("grpctls: server config doesn't use TLS")
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("grpctls: server TLS config provided is not valid: %v", err)
	}

	// without authentication
	if !cfg.ClientAuth {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("grpctls: could not load TLS keys: %s", err)
		}
		return creds, nil
	}

	// with authentication
	certificate, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("grpctls: could not load server key pair: %s", err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		return nil, fmt.Errorf("grpctls: could not read CA cert '%s': %v", cfg.CACert, err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("grpctls: configuring client's CA cert '%s'", cfg.CACert)
	}
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})
	return creds, nil
}

// Listener returns a valid listener server from an URI.
func Listener(uri string) (net.Listener, error) {
	proto, addr, err := ParseURI(uri)
	if err != nil {
		return nil, fmt.Errorf("grpctls: cannot parse address '%v': %v", uri, err)
	}
	lis, err := net.Listen(proto, addr)
	if err != nil {
		return nil, fmt.Errorf("grpctls: cannot listen socket '%v': %v", uri, err)
	}
	return lis, nil
}
