// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package resolvcheck

import "net"

// ClientMap manages mapping of client ips. It's not safe concurrency.
type ClientMap struct {
	m map[string]net.IP
}

// NewClientMap creates a new client map
func NewClientMap() *ClientMap {
	return &ClientMap{m: make(map[string]net.IP)}
}

// Set a new map between src and dst
func (c *ClientMap) Set(src, dst net.IP) {
	c.m[src.String()] = dst
}

// Get return mapped ip
func (c *ClientMap) Get(src net.IP) net.IP {
	dst, ok := c.m[src.String()]
	if !ok {
		return src
	}
	return dst
}
