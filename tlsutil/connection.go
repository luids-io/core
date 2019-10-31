// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsutil

import (
	"time"
)

// ConnectionData stores information from tls connections
type ConnectionData struct {
	ID   string          `json:"id" bson:"_id"`
	Info *ConnectionInfo `json:"info,omitempty" bson:",omitempty"`

	SendStream *StreamData `json:"sendStream,omitempty" bson:",omitempty"`
	RcvdStream *StreamData `json:"rcvdStream,omitempty" bson:",omitempty"`

	ClientHello *ClientHelloData `json:"clientHello,omitempty" bson:",omitempty"`
	ServerHello *ServerHelloData `json:"serverHello,omitempty" bson:",omitempty"`

	ClientCerts []CertSummary `json:"clientCerts,omitempty" bson:",omitempty"`
	ServerCerts []CertSummary `json:"serverCerts,omitempty" bson:",omitempty"`

	Tags []string `json:"tags,omitempty" bson:",omitempty"`
}

// ConnectionInfo stores main information from a tls connection
type ConnectionInfo struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`

	ClientIP   string `json:"clientIP"`
	ClientPort int    `json:"clientPort"`
	ServerIP   string `json:"serverIP"`
	ServerPort int    `json:"serverPort"`

	Uncompleted        bool `json:"uncompleted"`
	DetectedError      bool `json:"detectedError"`
	CompletedHandshake bool `json:"completedHandshake"`
}
