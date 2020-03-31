// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"time"

	"github.com/luids-io/core/tlsutil/layer"
)

// RecordData stores metadata of tls records
type RecordData struct {
	StreamID   string            `json:"streamID"`
	Timestamp  time.Time         `json:"timestamp"`
	Type       layer.ContentType `json:"type"`
	Len        uint16            `json:"len"`
	Ciphered   bool              `json:"ciphered"`
	Fragmented bool              `json:"fragmented" bson:",omitempty"`
	NumMsg     int               `json:"numMsg" bson:",omitempty"`
}
