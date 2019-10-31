// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package capture

import (
	"github.com/google/gopacket"
)

// Analyzer interface defines analyzer methods
type Analyzer interface {
	SendEtherPacket(gopacket.Packet) error
}

//Processor interface defines packet processor
type Processor interface {
	Register(key string, src *gopacket.PacketSource) error
	Unregister(key string) error
}

// Writer interface is used for write packages to persistant storage
type Writer interface {
	WritePacket(ci gopacket.CaptureInfo, data []byte) error
}
