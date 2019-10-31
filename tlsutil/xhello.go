// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"github.com/luisguillenc/tlslayer"
	"github.com/luisguillenc/tlslayer/tlsproto"
)

// ClientHelloData stores clienthello information
type ClientHelloData struct {
	ClientVersion   tlslayer.ProtocolVersion     `json:"clientVersion"`
	RandomLen       int                          `json:"randomLen"`
	SessionIDLen    int                          `json:"sessionIDLen"`
	SessionID       []byte                       `json:"sessionID,omitempty" bson:",omitempty"`
	CipherSuitesLen int                          `json:"cipherSuitesLen"`
	CipherSuites    []tlsproto.CipherSuite       `json:"cipherSuites"`
	CompressMethods []tlsproto.CompressionMethod `json:"compressMethods"`

	ExtensionLen  int             `json:"extensionLen"`
	Extensions    []ExtensionItem `json:"extensions,omitempty" bson:",omitempty"`
	ExtensionInfo *DecodedInfo    `json:"extensionInfo,omitempty" bson:",omitempty"`

	UseGREASE bool   `json:"useGREASE"`
	JA3       string `json:"ja3"`
	JA3digest string `json:"ja3digest"`
}

// ServerHelloData stores serverhello information
type ServerHelloData struct {
	ServerVersion     tlslayer.ProtocolVersion   `json:"serverVersion"`
	RandomLen         int                        `json:"randomLen"`
	SessionIDLen      int                        `json:"sessionIDLen"`
	SessionID         []byte                     `json:"sessionID,omitempty" bson:",omitempty"`
	CipherSuiteSel    tlsproto.CipherSuite       `json:"cipherSuiteSel"`
	CompressMethodSel tlsproto.CompressionMethod `json:"compressMethodSel"`

	ExtensionLen  int             `json:"extensionLen"`
	Extensions    []ExtensionItem `json:"extensions,omitempty" bson:",omitempty"`
	ExtensionInfo *DecodedInfo    `json:"extensionInfo,omitempty" bson:",omitempty"`
}

// ExtensionItem stores metadata information of extensions
type ExtensionItem struct {
	Type tlsproto.ExtensionType `json:"type"`
	Len  uint16                 `json:"len"`
}

// DecodedInfo stores information of extensions
type DecodedInfo struct {
	SNI                 string                        `json:"sni,omitempty" bson:",omitempty"`
	SignatureSchemes    []tlsproto.SignatureScheme    `json:"signatureSchemes,omitempty" bson:",omitempty"`
	SupportedVersions   []tlsproto.SupportedVersion   `json:"supportedVersions,omitempty" bson:",omitempty"`
	SupportedGroups     []tlsproto.SupportedGroup     `json:"supportedGroups,omitempty" bson:",omitempty"`
	ECPointFormats      []tlsproto.ECPointFormat      `json:"ecPointFormats,omitempty" bson:",omitempty"`
	OSCP                bool                          `json:"oscp"`
	ALPNs               []string                      `json:"alpns,omitempty" bson:",omitempty"`
	KeyShareEntries     []tlsproto.KeyShareEntry      `json:"keyShareEntries,omitempty" bson:",omitempty"`
	PSKKeyExchangeModes []tlsproto.PSKKeyExchangeMode `json:"pskKeyExchangeModes,omitempty" bson:",omitempty"`
}
