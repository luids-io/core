// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package grpctls

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ClientCfg defines configuration for a client.
type ClientCfg struct {
	// CertFile path to the client certificate
	CertFile string `json:"certfile,omitempty"`
	// KeyFile path to the private key
	KeyFile string `json:"keyfile,omitempty"`
	// ServerName is used for server check validation
	ServerName string `json:"servername,omitempty"`
	// ServerCert path to the server certificate
	ServerCert string `json:"servercert,omitempty"`
	// CACert path to certification authority certificate
	CACert string `json:"cacert,omitempty"`
	// UseSystemCAs if client uses system wide CA certs
	UseSystemCAs bool `json:"systemca"`
}

// UseTLS returns true if TLS configuration is set.
func (cfg ClientCfg) UseTLS() bool {
	if cfg.ServerCert == "" && cfg.CACert == "" && !cfg.UseSystemCAs {
		return false
	}
	return true
}

// Empty returns true if configuration values are empty.
func (cfg ClientCfg) Empty() bool {
	if cfg.CertFile != "" {
		return false
	}
	if cfg.KeyFile != "" {
		return false
	}
	if cfg.ServerName != "" {
		return false
	}
	if cfg.ServerCert != "" {
		return false
	}
	if cfg.CACert != "" {
		return false
	}
	if cfg.UseSystemCAs {
		return false
	}
	return true
}

// Validate if configuration is ok.
func (cfg ClientCfg) Validate() error {
	if cfg.CertFile != "" {
		if !fileExists(cfg.CertFile) {
			return fmt.Errorf("certfile '%v' doesn't exists", cfg.CertFile)
		}
		if cfg.KeyFile == "" {
			return errors.New("keyfile is required for client's certificate")
		}
		if !fileExists(cfg.KeyFile) {
			return fmt.Errorf("keyfile '%v' doesn't exists", cfg.KeyFile)
		}
	}
	if cfg.ServerCert != "" {
		if !fileExists(cfg.ServerCert) {
			return fmt.Errorf("servercert file '%v' doesn't exists", cfg.ServerCert)
		}
	}
	if cfg.CACert != "" {
		if !fileExists(cfg.CACert) {
			return fmt.Errorf("cacert file '%v' doesn't exists", cfg.CACert)
		}
	}
	//some logic
	if cfg.ServerCert != "" && (cfg.CACert != "" || cfg.UseSystemCAs) {
		return fmt.Errorf("using a server certfile and CA configuration")
	}
	if cfg.CertFile != "" && (cfg.CACert == "" && !cfg.UseSystemCAs) {
		return fmt.Errorf("CA configuration is needed for client side auth")
	}
	return nil
}

func (cfg ClientCfg) getCreds(addr string) (credentials.TransportCredentials, error) {
	if cfg.ServerCert != "" {
		return credentials.NewClientTLSFromFile(cfg.ServerCert, "")
	}
	var err error
	var certPool *x509.CertPool
	if cfg.UseSystemCAs {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("can't get system cert pool: %v", err)
		}
	} else {
		certPool = x509.NewCertPool()
	}
	if cfg.CACert != "" {
		ca, err := ioutil.ReadFile(cfg.CACert)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate '%s': %v", cfg.CACert, err)
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append client certs")
		}
	}
	tlsConfig := &tls.Config{RootCAs: certPool}
	if cfg.ServerName != "" {
		tlsConfig.ServerName = cfg.ServerName
	} else {
		servername, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("could not get servername from '%s': %v", addr, err)
		}
		tlsConfig.ServerName = servername
	}
	if cfg.CertFile != "" {
		certificate, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{certificate}
	}
	return credentials.NewTLS(tlsConfig), nil
}

// Dial is used for grpc client dialing
func Dial(uri string, cfg ClientCfg, grpcOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), uri, cfg, grpcOpts...)
}

// DialContext is used for grpc client dialing with context
func DialContext(ctx context.Context, uri string, cfg ClientCfg, grpcOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	proto, addr, err := ParseURI(uri)
	if err != nil {
		return nil, fmt.Errorf("grpctls: cannot parse URI '%v': %v", uri, err)
	}
	if proto == "unix" {
		dopts := make([]grpc.DialOption, 0)
		dopts = append(dopts, grpc.WithInsecure())
		dopts = append(dopts, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
		dopts = append(dopts, grpcOpts...)
		return grpc.DialContext(ctx, addr, dopts...)
	}
	//proto == tcp
	if !cfg.UseTLS() {
		dopts := make([]grpc.DialOption, 0)
		dopts = append(dopts, grpc.WithInsecure())
		dopts = append(dopts, grpcOpts...)
		return grpc.DialContext(ctx, addr, dopts...)
	}
	//useTLS
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("grpctls: validating client tls config: %v", err)
	}
	creds, err := cfg.getCreds(addr)
	if err != nil {
		return nil, fmt.Errorf("grpctls: getting client tls credentials: %v", err)
	}
	dopts := make([]grpc.DialOption, 0)
	dopts = append(dopts, grpc.WithTransportCredentials(creds))
	dopts = append(dopts, grpcOpts...)
	return grpc.DialContext(ctx, addr, dopts...)
}
