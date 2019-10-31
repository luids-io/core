// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"time"

	"github.com/luisguillenc/tlslayer"
)

// RecordData stores metadata of tls records
type RecordData struct {
	StreamID   string               `json:"streamID"`
	Timestamp  time.Time            `json:"timestamp"`
	Type       tlslayer.ContentType `json:"type"`
	Len        uint16               `json:"len"`
	Ciphered   bool                 `json:"ciphered"`
	Fragmented bool                 `json:"fragmented" bson:",omitempty"`
	NumMsg     int                  `json:"numMsg" bson:",omitempty"`
}
