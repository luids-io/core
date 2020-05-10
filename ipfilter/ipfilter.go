// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.md.

// Package ipfilter provides a simple package for checking GRPC/HTTP requests
// against a white/black list of IPs or CIDRs.
//
// This package is a work in progress and makes no API stability promises.
package ipfilter

import (
	"net"
	"net/http"
)

// Action to be taken with the request
type Action bool

// Posible values for Action
const (
	Deny   Action = false
	Accept Action = true
)

// Filter is the main structure
type Filter struct {
	// Wrapped http handler to control
	Wrapped http.Handler

	AllowedIPs      []net.IP
	DisallowedIPs   []net.IP
	AllowedCIDRs    []*net.IPNet
	DisallowedCIDRs []*net.IPNet
	// Policy is the action to take if no match
	Policy Action
}

// Empty resturns true if no list has information
func (f Filter) Empty() bool {
	if f.AllowedIPs != nil && len(f.AllowedIPs) > 0 {
		return false
	}
	if f.DisallowedIPs != nil && len(f.DisallowedIPs) > 0 {
		return false
	}
	if f.AllowedCIDRs != nil && len(f.AllowedCIDRs) > 0 {
		return false
	}
	if f.DisallowedCIDRs != nil && len(f.DisallowedCIDRs) > 0 {
		return false
	}
	return true
}

// Check returns the action to be taken for the ip
func (f Filter) Check(ip net.IP) Action {
	if f.AllowedCIDRs != nil && len(f.AllowedCIDRs) > 0 {
		for _, net := range f.AllowedCIDRs {
			if net.Contains(ip) {
				return Accept
			}
		}
	}
	if f.AllowedIPs != nil && len(f.AllowedIPs) > 0 {
		for _, allowed := range f.AllowedIPs {
			if ip.Equal(allowed) {
				return Accept
			}
		}
	}
	if f.DisallowedCIDRs != nil && len(f.DisallowedCIDRs) > 0 {
		for _, net := range f.DisallowedCIDRs {
			if net.Contains(ip) {
				return Deny
			}
		}
	}
	if f.DisallowedIPs != nil && len(f.DisallowedIPs) > 0 {
		for _, allowed := range f.DisallowedIPs {
			if ip.Equal(allowed) {
				return Deny
			}
		}
	}
	return f.Policy
}

// Whitelist creates a filter with the list of ips and/or cidrs passed.
// Returned filter has a deny policy and all ips and/or cidrs passed
// will be added to the allowed lists.
// Strings that can't be parsed will be ignored.
func Whitelist(allowed []string) Filter {
	filter := Filter{}
	filter.Policy = Deny
	filter.AllowedIPs = make([]net.IP, 0)
	filter.AllowedCIDRs = make([]*net.IPNet, 0)
	for _, item := range allowed {
		_, cidr, err := net.ParseCIDR(item)
		if err == nil {
			filter.AllowedCIDRs = append(filter.AllowedCIDRs, cidr)
			continue
		}
		ip := net.ParseIP(item)
		if ip != nil {
			filter.AllowedIPs = append(filter.AllowedIPs, ip)
		}
	}
	return filter
}

// Blacklist creates a filter with the list of ips and/or cidrs passed.
// Returned filter has an allow policy and all ips and/or cidrs passed
// will be added to the disallowed lists.
// Strings that can't be parsed will be ignored.
func Blacklist(disallowed []string) Filter {
	filter := Filter{}
	filter.Policy = Accept
	filter.DisallowedIPs = make([]net.IP, 0)
	filter.DisallowedCIDRs = make([]*net.IPNet, 0)
	for _, item := range disallowed {
		_, cidr, err := net.ParseCIDR(item)
		if err == nil {
			filter.DisallowedCIDRs = append(filter.DisallowedCIDRs, cidr)
			continue
		}
		ip := net.ParseIP(item)
		if ip != nil {
			filter.DisallowedIPs = append(filter.DisallowedIPs, ip)
		}
	}
	return filter
}
