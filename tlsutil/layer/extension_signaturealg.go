// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import "fmt"

// SignatureScheme defines a signature scheme defined in https://www.iana.org/assignments/tls-parameters/tls-parameters.xhtml#tls-signaturescheme
type SignatureScheme uint16

// IsGREASE returns true if passed signature and hash is in GREASE spec
func (s SignatureScheme) IsGREASE() bool {
	return isGREASE16(uint16(s))
}

func (s SignatureScheme) getDesc() string {
	if s.IsGREASE() {
		return "GREASE"
	}
	n := uint16(s)
	switch n {
	case 0x0201:
		return "rsa_pkcs1_sha1"
	case 0x0203:
		return "ecdsa_sha1"
	case 0x0401:
		return "rsa_pkcs1_sha256"
	case 0x0403:
		return "ecdsa_secp256r1_sha256"
	case 0x0501:
		return "rsa_pkcs1_sha384"
	case 0x0503:
		return "ecdsa_secp384r1_sha384"
	case 0x0601:
		return "rsa_pkcs1_sha512"
	case 0x0603:
		return "ecdsa_secp521r1_sha512"
	case 0x0804:
		return "rsa_pss_rsae_sha256"
	case 0x0805:
		return "rsa_pss_rsae_sha384"
	case 0x0806:
		return "rsa_pss_rsae_sha512"
	case 0x0807:
		return "ed25519"
	case 0x0808:
		return "ed448"
	case 0x0809:
		return "rsa_pss_pss_sha256"
	case 0x080A:
		return "rsa_pss_pss_sha384"
	case 0x080B:
		return "rsa_pss_pss_sha512"
	}
	return "unknown"
}

func (s SignatureScheme) String() string {
	return fmt.Sprintf("%s(%d)", s.getDesc(), s)
}
