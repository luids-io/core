// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

// HandshakeType defines the type of handshake
type HandshakeType uint8

// Constants of HandshakeType
const (
	HandshakeTypeHelloRequest       HandshakeType = 0
	HandshakeTypeClientHello        HandshakeType = 1
	HandshakeTypeServerHello        HandshakeType = 2
	HandshakeTypeNewSessionTicket   HandshakeType = 4
	HandshakeTypeEndOfEarlyData     HandshakeType = 5
	HandshakeTypeCertificate        HandshakeType = 11
	HandshakeTypeServerKeyExchange  HandshakeType = 12
	HandshakeTypeCertificateRequest HandshakeType = 13
	HandshakeTypeServerHelloDone    HandshakeType = 14
	HandshakeTypeCertificateVerify  HandshakeType = 15
	HandshakeTypeClientKeyExchange  HandshakeType = 16
	HandshakeTypeFinished           HandshakeType = 20
	HandshakeTypeCertificateURL     HandshakeType = 21
	HandshakeTypeCertificateStatus  HandshakeType = 22
	HandshakeTypeKeyUpdate          HandshakeType = 24
)

// HandShakeTypeReg is a map with strings of alert description
var handShakeTypeReg = map[HandshakeType]string{
	HandshakeTypeHelloRequest:       "hello_request",
	HandshakeTypeClientHello:        "client_hello",
	HandshakeTypeServerHello:        "server_hello",
	HandshakeTypeNewSessionTicket:   "new_session_ticket",
	HandshakeTypeEndOfEarlyData:     "end_of_early_data",
	HandshakeTypeCertificate:        "certificate",
	HandshakeTypeServerKeyExchange:  "server_key_exchange",
	HandshakeTypeCertificateRequest: "certificate_request",
	HandshakeTypeServerHelloDone:    "server_hello_done",
	HandshakeTypeCertificateVerify:  "certificate_verify",
	HandshakeTypeClientKeyExchange:  "client_key_exchange",
	HandshakeTypeFinished:           "finished",
	HandshakeTypeCertificateURL:     "certificate_url",
	HandshakeTypeCertificateStatus:  "certificate_status",
	HandshakeTypeKeyUpdate:          "key_update",
}

func (hst HandshakeType) getDesc() string {
	if h, ok := handShakeTypeReg[hst]; ok {
		return h
	}
	return "unknown"
}

func (hst HandshakeType) String() string {
	return fmt.Sprintf("%s(%d)", hst.getDesc(), hst)
}

// IsValid method checks if it's a valid value
func (hst HandshakeType) IsValid() bool {
	_, ok := handShakeTypeReg[hst]
	return ok
}
