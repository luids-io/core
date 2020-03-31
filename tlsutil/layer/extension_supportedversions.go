// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package layer

import (
	"fmt"
)

// SupportedVersion is like ProtocolVersion but values "0x7f00 | draft_version" are valid
type SupportedVersion uint16

// IsGREASE returns true if is a grease value
func (sv SupportedVersion) IsGREASE() bool {
	return isGREASE16(uint16(sv))
}

// IsDraft returns true if is a draft version
func (sv SupportedVersion) IsDraft() bool {
	mask := SupportedVersion(0xff00)
	r := mask & sv
	return r == SupportedVersion(0x7f00)
}

func (sv SupportedVersion) getDesc() string {
	if sv.IsDraft() {
		ver := SupportedVersion(0x00FF) & sv
		return fmt.Sprintf("TLS_1.3(draft %d)", ver)
	}
	if sv.IsGREASE() {
		return "GREASE"
	}
	switch ProtocolVersion(sv) {
	case VersionSSL30:
		return "SSL_3.0"
	case VersionTLS10:
		return "TLS_1.0"
	case VersionTLS11:
		return "TLS_1.1"
	case VersionTLS12:
		return "TLS_1.2"
	case VersionTLS13:
		return "TLS_1.3"
	default:
		return "unknown"
	}
}

func (sv SupportedVersion) String() string {
	return fmt.Sprintf("%s(%d)", sv.getDesc(), sv)
}
