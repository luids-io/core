// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package packetproc

import (
	"github.com/google/gopacket"
)

// Analyzer interface defines analyzer methods
type Analyzer interface {
	SendEtherPacket(data []byte, md *gopacket.PacketMetadata) error
	SendIPv4Packet(data []byte, md *gopacket.PacketMetadata) error
	SendIPv6Packet(data []byte, md *gopacket.PacketMetadata) error
}

// Writer interface is used for write packages to persistant storage
type Writer interface {
	WritePacket(ci gopacket.CaptureInfo, data []byte) error
}
