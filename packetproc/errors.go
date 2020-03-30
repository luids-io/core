// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package packetproc

import (
	"fmt"

	"github.com/google/gopacket"
)

// ErrorsBuffer sets the default size for error channels
var ErrorsBuffer = 20

// ShowPacketInError if a packet digest will be show with the error string
var ShowPacketInError = false

// Error is used for packet processing
type Error struct {
	packet gopacket.Packet
	err    error
}

// NewError creates a new packet processing error
func NewError(packet gopacket.Packet, err error) *Error {
	return &Error{packet: packet, err: err}
}

// Error implements error interface
func (e *Error) Error() string {
	serr := e.err.Error()
	if ShowPacketInError {
		return fmt.Sprintf("%s [%v]", serr, e.packet)
	}
	return serr
}

func (e *Error) String() string {
	return e.Error()
}

// Internal returns internal error
func (e *Error) Internal() error {
	return e.err
}

// Packet returns packet with error
func (e *Error) Packet() gopacket.Packet {
	return e.packet
}
