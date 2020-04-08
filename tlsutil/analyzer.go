// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"net"
	"time"
)

// AnalyzerFactory interface is used for create Analyzer services
type AnalyzerFactory interface {
	NewAnalyzer() Analyzer
}

// Analyzer interface defines analyzer methods
type Analyzer interface {
	SendMessage(m *Msg) error
	Close() error
}

// Msg defines message for analyzer
type Msg struct {
	Type     MsgType
	StreamID int64
	Open     *MsgOpen
	Data     *MsgData
}

// MsgType defines message types
type MsgType int8

// Type possible values
const (
	DataMsg MsgType = iota
	OpenMsg
	CloseMsg
)

// MsgOpen stores required data by the open message
type MsgOpen struct {
	SrcIP, DstIP     net.IP
	SrcPort, DstPort int
}

// MsgData stores required data by the data message
type MsgData struct {
	Timestamp        time.Time
	Bytes            int
	SawStart, SawEnd bool
	Records          [][]byte
	Error            error
}
