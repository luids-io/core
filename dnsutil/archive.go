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

// ResolvData defines struct for archive
type ResolvData struct {
	ID        string    `json:"id" bson:"_id"`
	Timestamp time.Time `json:"timestamp"`
	Server    net.IP    `json:"server"`
	Client    net.IP    `json:"client"`
	Resolved  []net.IP  `json:"resolved"`
	Name      string    `json:"name"`
}
