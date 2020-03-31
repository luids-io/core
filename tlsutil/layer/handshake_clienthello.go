// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

const (
	clientHelloRandomLen = 32
)

// ClientHelloData stores data from a clienthello handshake
type ClientHelloData struct {
	ClientVersion   ProtocolVersion     `json:"clientVersion"`
	Random          []byte              `json:"random,omitempty"`
	SessionID       []byte              `json:"sessionID,omitempty"`
	CipherSuites    []CipherSuite       `json:"cipherSuites,omitempty"`
	CompressMethods []CompressionMethod `json:"compressMethods,omitempty"`

	ExtensionsLen uint16          `json:"extensionsLen"`
	Extensions    []Extension     `json:"extensions,omitempty"`
	ExtInfo       *ExtensionsInfo `json:"extInfo,omitempty"`
}

func (ch *ClientHelloData) String() string {
	str := fmt.Sprintln("Version:", ch.ClientVersion)
	str += fmt.Sprintf("SessionID: %#v\n", ch.SessionID)
	str += fmt.Sprintf("Cipher Suites: %v\n", ch.CipherSuites)
	str += fmt.Sprintf("Compression Methods: %v\n", ch.CompressMethods)
	str += fmt.Sprintf("Extensions: %v\n", ch.Extensions)
	str += fmt.Sprintln("Extensions info:", ch.ExtInfo)

	return str
}

// UseGREASE returns true if clienthello data uses GREASE proposal
func (ch *ClientHelloData) UseGREASE() bool {
	for _, c := range ch.CipherSuites {
		if c.IsGREASE() {
			return true
		}
	}
	for _, e := range ch.Extensions {
		if e.Type.IsGREASE() {
			return true
		}
	}
	if ch.ExtInfo != nil {
		if len(ch.ExtInfo.SignatureSchemes) > 0 {
			for _, sa := range ch.ExtInfo.SignatureSchemes {
				if sa.IsGREASE() {
					return true
				}
			}
		}
		if len(ch.ExtInfo.SupportedGroups) > 0 {
			for _, sg := range ch.ExtInfo.SupportedGroups {
				if sg.IsGREASE() {
					return true
				}
			}
		}
		if len(ch.ExtInfo.SupportedVersions) > 0 {
			for _, sv := range ch.ExtInfo.SupportedVersions {
				if sv.IsGREASE() {
					return true
				}
			}
		}
	}
	return false
}
