// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

// Extension stores extension information
type Extension struct {
	Type    ExtensionType `json:"type"`
	Len     uint16        `json:"len"`
	payload []byte
}

// ExtensionsInfo stores all decoded information from extensions
type ExtensionsInfo struct {
	// ExtServerName
	SNI string `json:"sni,omitempty"`
	// ExtSignatureAlgs
	SignatureSchemes []SignatureScheme `json:"signatureSchemes,omitempty"`
	// ExtSupportedVersions
	SupportedVersions []SupportedVersion `json:"supportedVersions,omitempty"`
	// ExtSupportedGroups
	SupportedGroups []SupportedGroup `json:"supportedGroups,omitempty"`
	// ExtECPointFormats
	ECPointFormats []ECPointFormat `json:"ecPointFormats,omitempty"`
	// ExtStatusRequest
	OSCP bool `json:"oscp"`
	// ExtALPN
	ALPNs []string `json:"alpns,omitempty"`
	// ExtKeyShare
	KeyShareEntries []KeyShareEntry `json:"keyShareEntries,omitempty"`
	// ExtPSKKeyExchangeModes
	PSKKeyExchangeModes []PSKKeyExchangeMode `json:"pskKeyExchangeModes,omitempty"`
}

// ExtensionType is an extension type defined by rfc
type ExtensionType uint16

// TLS Extensions http://www.iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml
const (
	ExtServerName           ExtensionType = 0
	ExtMaxFragLen           ExtensionType = 1
	ExtClientCertURL        ExtensionType = 2
	ExtTrustedCAKeys        ExtensionType = 3
	ExtTruncatedHMAC        ExtensionType = 4
	ExtStatusRequest        ExtensionType = 5
	ExtUserMapping          ExtensionType = 6
	ExtClientAuthz          ExtensionType = 7
	ExtServerAuthz          ExtensionType = 8
	ExtCertType             ExtensionType = 9
	ExtSupportedGroups      ExtensionType = 10
	ExtECPointFormats       ExtensionType = 11
	ExtSRP                  ExtensionType = 12
	ExtSignatureAlgs        ExtensionType = 13
	ExtUseSRTP              ExtensionType = 14
	ExtHeartbeat            ExtensionType = 15
	ExtALPN                 ExtensionType = 16 // Replaced NPN
	ExtStatusRequestV2      ExtensionType = 17
	ExtSignedCertTS         ExtensionType = 18 // Certificate Transparency
	ExtClientCertType       ExtensionType = 19
	ExtServerCertType       ExtensionType = 20
	ExtPadding              ExtensionType = 21 // Temp http://www.iana.org/go/draft-ietf-tls-padding
	ExtEncryptThenMAC       ExtensionType = 22
	ExtExtendedMasterSecret ExtensionType = 23
	ExtTokenBinding         ExtensionType = 24
	ExtCachedInfo           ExtensionType = 25
	ExtCompressCert         ExtensionType = 27
	ExtRecordSizeLimit      ExtensionType = 28
	ExtPwdProtect           ExtensionType = 29
	ExtPwdClear             ExtensionType = 30
	ExtPasswordSalt         ExtensionType = 31
	ExtSessionTicket        ExtensionType = 35
	ExtPreSharedKey         ExtensionType = 41
	ExtEarlyData            ExtensionType = 42
	ExtSupportedVersions    ExtensionType = 43
	ExtCookie               ExtensionType = 44
	ExtPSKKeyExchangeModes  ExtensionType = 45
	ExtCertAuthorities      ExtensionType = 47
	ExtOIDFilters           ExtensionType = 48
	ExtPostHandshakeAuth    ExtensionType = 49
	ExtSignatureAlgsCert    ExtensionType = 50
	ExtKeyShare             ExtensionType = 51
	ExtNPN                  ExtensionType = 13172 // Next Protocol Negotiation not ratified and replaced by ALPN
	ExtRenegotiationInfo    ExtensionType = 65281
)

var extensionReg = map[ExtensionType]string{
	ExtServerName:           "server_name",
	ExtMaxFragLen:           "max_fragment_length",
	ExtClientCertURL:        "client_certificate_url",
	ExtTrustedCAKeys:        "trusted_ca_keys",
	ExtTruncatedHMAC:        "truncated_hmac",
	ExtStatusRequest:        "status_request",
	ExtUserMapping:          "user_mapping",
	ExtClientAuthz:          "client_authz",
	ExtServerAuthz:          "server_authz",
	ExtCertType:             "cert_type",
	ExtSupportedGroups:      "supported_groups",
	ExtECPointFormats:       "ec_point_formats",
	ExtSRP:                  "srp",
	ExtSignatureAlgs:        "signature_algorithms",
	ExtUseSRTP:              "use_srtp",
	ExtHeartbeat:            "heartbeat",
	ExtALPN:                 "application_layer_protocol_negotiation",
	ExtStatusRequestV2:      "status_request_v2",
	ExtSignedCertTS:         "signed_certificate_timestamp",
	ExtClientCertType:       "client_certificate_type",
	ExtServerCertType:       "server_certificate_type",
	ExtPadding:              "padding",
	ExtEncryptThenMAC:       "encrypt_then_mac",
	ExtExtendedMasterSecret: "extended_master_secret",
	ExtTokenBinding:         "token_binding",
	ExtCachedInfo:           "cached_info",
	ExtCompressCert:         "compress_certificate",
	ExtRecordSizeLimit:      "record_size_limit",
	ExtPwdProtect:           "pwd_protect",
	ExtPwdClear:             "pwd_clear",
	ExtPasswordSalt:         "password_salt",
	ExtSessionTicket:        "session_ticket",
	ExtPreSharedKey:         "pre_shared_key",
	ExtEarlyData:            "early_data",
	ExtSupportedVersions:    "supported_versions",
	ExtCookie:               "cookie",
	ExtPSKKeyExchangeModes:  "psk_key_exchange_modes",
	ExtCertAuthorities:      "certificate_authorities",
	ExtOIDFilters:           "oid_filters",
	ExtPostHandshakeAuth:    "post_handshake_auth",
	ExtSignatureAlgsCert:    "signature_algorithms_cert",
	ExtKeyShare:             "key_share",
	ExtNPN:                  "next_protocol_negotiation",
	ExtRenegotiationInfo:    "renegotiation_info",
}

func (e ExtensionType) getDesc() string {
	if ext, ok := extensionReg[e]; ok {
		return ext
	}
	if e.IsGREASE() {
		return "GREASE"
	}
	return "unknown"
}

// String method for a TLS Extension
func (e ExtensionType) String() string {
	return fmt.Sprintf("%s(%d)", e.getDesc(), e)
}

// IsGREASE returns true if is an extension reserved by GREASE rfc
func (e ExtensionType) IsGREASE() bool {
	return isGREASE16(uint16(e))
}

func (e Extension) String() string {
	return fmt.Sprintf("%s (len=%d)", e.Type, e.Len)
}

func (i *ExtensionsInfo) String() string {
	str := fmt.Sprintf("SNI: %q\n", i.SNI)
	str += fmt.Sprintf("Signature Schemes: %v\n", i.SignatureSchemes)
	str += fmt.Sprintf("Supported Groups: %v\n", i.SupportedGroups)
	str += fmt.Sprintf("ECPoints Formats: %v\n", i.ECPointFormats)
	str += fmt.Sprintf("OSCP: %v\n", i.OSCP)
	str += fmt.Sprintf("ALPNs: %v", i.ALPNs)
	str += fmt.Sprintf("Supported Versions: %v\n", i.SupportedVersions)
	str += fmt.Sprintf("Key Share Entries: %v\n", i.KeyShareEntries)
	str += fmt.Sprintf("PSK Key Exchange Modes: %v\n", i.PSKKeyExchangeModes)

	return str
}
