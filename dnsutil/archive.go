// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package dnsutil

import (
	"context"
	"net"
	"time"
)

// Archiver is the interface for archiver
type Archiver interface {
	SaveResolv(context.Context, ResolvData) (string, error)
}

// ResolvData defines struct for summary ip resolutions
type ResolvData struct {
	ID        string        `json:"id" bson:"_id"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	Server    net.IP        `json:"server"`
	Client    net.IP        `json:"client"`
	//query info
	QID              uint16 `json:"qid"`
	Name             string `json:"name"`
	CheckingDisabled bool   `json:"checkingDisabled"`
	//response info
	ReturnCode        uint16   `json:"returnCode"`
	AuthenticatedData bool     `json:"authenticatedData"`
	Resolved          []net.IP `json:"resolved,omitempty" bson:",omitempty"`
}
