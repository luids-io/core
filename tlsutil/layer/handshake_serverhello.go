// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

const (
	serverHelloRandomLen = 32
)

// ServerHelloData stores data from ServerHello messages
type ServerHelloData struct {
	ServerVersion     ProtocolVersion   `json:"serverVersion"`
	Random            []byte            `json:"random,omitempty"`
	SessionID         []byte            `json:"sessionID,omitempty"`
	CipherSuiteSel    CipherSuite       `json:"cipherSuiteSel"`
	CompressMethodSel CompressionMethod `json:"compressMethodSel"`

	ExtensionsLen uint16          `json:"extensionsLen"`
	Extensions    []Extension     `json:"extensions,omitempty"`
	ExtInfo       *ExtensionsInfo `json:"extInfo,omitempty"`
}

func (hs *ServerHelloData) String() string {
	str := fmt.Sprintln("Version:", hs.ServerVersion)
	str += fmt.Sprintf("SessionID: %#v\n", hs.SessionID)
	str += fmt.Sprintf("Cipher Suite selected: %v\n", hs.CipherSuiteSel)
	str += fmt.Sprintf("Compression selected: %v\n", hs.CompressMethodSel)
	str += fmt.Sprintln("Extensions:", hs.Extensions)
	str += fmt.Sprintln("Extensions info:", hs.ExtInfo)

	return str
}
