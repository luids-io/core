// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"github.com/luids-io/core/tlsutil/layer"
)

// ClientHelloData stores clienthello information
type ClientHelloData struct {
	ClientVersion   layer.ProtocolVersion     `json:"clientVersion"`
	RandomLen       int                       `json:"randomLen"`
	SessionIDLen    int                       `json:"sessionIDLen"`
	SessionID       []byte                    `json:"sessionID,omitempty" bson:",omitempty"`
	CipherSuitesLen int                       `json:"cipherSuitesLen"`
	CipherSuites    []layer.CipherSuite       `json:"cipherSuites"`
	CompressMethods []layer.CompressionMethod `json:"compressMethods"`

	ExtensionLen  int             `json:"extensionLen"`
	Extensions    []ExtensionItem `json:"extensions,omitempty" bson:",omitempty"`
	ExtensionInfo *DecodedInfo    `json:"extensionInfo,omitempty" bson:",omitempty"`

	UseGREASE bool   `json:"useGREASE"`
	JA3       string `json:"ja3"`
	JA3digest string `json:"ja3digest"`
}

// ServerHelloData stores serverhello information
type ServerHelloData struct {
	ServerVersion     layer.ProtocolVersion   `json:"serverVersion"`
	RandomLen         int                     `json:"randomLen"`
	SessionIDLen      int                     `json:"sessionIDLen"`
	SessionID         []byte                  `json:"sessionID,omitempty" bson:",omitempty"`
	CipherSuiteSel    layer.CipherSuite       `json:"cipherSuiteSel"`
	CompressMethodSel layer.CompressionMethod `json:"compressMethodSel"`

	ExtensionLen  int             `json:"extensionLen"`
	Extensions    []ExtensionItem `json:"extensions,omitempty" bson:",omitempty"`
	ExtensionInfo *DecodedInfo    `json:"extensionInfo,omitempty" bson:",omitempty"`
}

// ExtensionItem stores metadata information of extensions
type ExtensionItem struct {
	Type layer.ExtensionType `json:"type"`
	Len  uint16              `json:"len"`
}

// DecodedInfo stores information of extensions
type DecodedInfo struct {
	SNI                 string                     `json:"sni,omitempty" bson:",omitempty"`
	SignatureSchemes    []layer.SignatureScheme    `json:"signatureSchemes,omitempty" bson:",omitempty"`
	SupportedVersions   []layer.SupportedVersion   `json:"supportedVersions,omitempty" bson:",omitempty"`
	SupportedGroups     []layer.SupportedGroup     `json:"supportedGroups,omitempty" bson:",omitempty"`
	ECPointFormats      []layer.ECPointFormat      `json:"ecPointFormats,omitempty" bson:",omitempty"`
	OSCP                bool                       `json:"oscp"`
	ALPNs               []string                   `json:"alpns,omitempty" bson:",omitempty"`
	KeyShareEntries     []layer.KeyShareEntry      `json:"keyShareEntries,omitempty" bson:",omitempty"`
	PSKKeyExchangeModes []layer.PSKKeyExchangeMode `json:"pskKeyExchangeModes,omitempty" bson:",omitempty"`
}
