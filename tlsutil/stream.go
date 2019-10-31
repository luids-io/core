// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"time"

	"github.com/luisguillenc/tlslayer/tlsproto"
)

// StreamData stores data of tls streams
type StreamData struct {
	ID   string      `json:"id" bson:"_id"`
	Info *StreamInfo `json:"info"`

	PlaintextAcc  *PlaintextSummary  `json:"plaintextAcc,omitempty" bson:",omitempty"`
	CiphertextAcc *CiphertextSummary `json:"ciphertextAcc,omitempty" bson:",omitempty"`

	HandshakeSeq []HandshakeItem `json:"handshakeSeq,omitempty" bson:",omitempty"`
	HandshakeSum int             `json:"handshakeSum,omitempty" bson:",omitempty"`
}

// StreamInfo stores information of stream
type StreamInfo struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
	SawStart bool          `json:"sawStart"`
	SawEnd   bool          `json:"sawEnd"`

	DetectedError bool      `json:"detectedError"`
	ErrorType     string    `json:"errorType,omitempty" bson:",omitempty"`
	ErrorTime     time.Time `json:"errorTime,omitempty" bson:",omitempty"`

	SrcIP4  string `json:"srcIP"`
	DstIP4  string `json:"dstIP"`
	SrcPort int    `json:"srcPort"`
	DstPort int    `json:"dstPort"`

	Bytes   int64   `json:"bytes"`
	Packets int64   `json:"packets"`
	BPS     float32 `json:"bps"`
	PPS     float32 `json:"pps"`
}

// HandshakeItem stores handshake metadata information
type HandshakeItem struct {
	Type tlsproto.HandshakeType `json:"type"`
	Len  uint32                 `json:"len"`
}

// PlaintextSummary stores summary of plaintext traffic
type PlaintextSummary struct {
	HskRecords          int64 `json:"hskRecords"`
	HskBytes            int64 `json:"hskBytes"`
	AlertRecords        int64 `json:"alertRecords"`
	AlertBytes          int64 `json:"alertBytes"`
	CCTRecords          int64 `json:"cctRecords"`
	CCTBytes            int64 `json:"cctBytes"`
	AppDataRecords      int64 `json:"appDataRecords"`
	AppDataBytes        int64 `json:"appDataBytes"`
	FragmentedRecords   int   `json:"fragmentedRecords"`
	MaxMessagesInRecord int   `json:"maxMessagesInRecord"`
}

// CiphertextSummary stores summary of ciphertext traffic
type CiphertextSummary struct {
	HskRecords     int64 `json:"hskRecords"`
	HskBytes       int64 `json:"hskBytes"`
	AlertRecords   int64 `json:"alertRecords"`
	AlertBytes     int64 `json:"alertBytes"`
	CCTRecords     int64 `json:"cctRecords"`
	CCTBytes       int64 `json:"cctBytes"`
	AppDataRecords int64 `json:"appDataRecords"`
	AppDataBytes   int64 `json:"appDataBytes"`
}
