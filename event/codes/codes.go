// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package codes includes "official" event codes and registers them.
//
// This package is a work in progress and makes no API stability promises.
package codes

import (
	"github.com/luids-io/core/event"
)

//A big enum with registered codes
const (
	//Security type codes
	TestSecurity event.Code = 10000

	//DNS Blackhole
	DNSListedDomain   event.Code = 10001
	DNSUnlistedDomain event.Code = 10002
	DNSListedIP       event.Code = 10003
	DNSUnlistedIP     event.Code = 10004

	//DNS Collect
	DNSMaxClientRequests  event.Code = 10005
	DNSMaxNamesResolvedIP event.Code = 10006

	//Netfiler queue
	NFQListedIP   event.Code = 10010
	NFQUnlistedIP event.Code = 10011

	//TLS processor
	TLSListedSNI    event.Code = 10020
	TLSUnlistedSNI  event.Code = 10021
	TLSInvalidCerts event.Code = 10022
	TLSNonTLSData   event.Code = 10023
)

var registry = []item{
	{TestSecurity, "TestSecurity", "Test event with data: [data.test]"},

	//DNS Blackhole
	{DNSListedDomain, "DNSListedDomain", "Domain '[data.listed]' listed has been resolved by [data.remote]"},
	{DNSUnlistedDomain, "DNSUnlistedDomain", "Domain '[data.listed]' unlisted has been resolved by [data.remote]"},
	{DNSListedIP, "DNSListedIP", "IP [data.listed] listed has been resolved by [data.remote]"},
	{DNSUnlistedIP, "DNSUnlistedIP", "IP [data.listed] unlisted has been resolved by [data.remote]"},

	//DNS Collect
	{DNSMaxClientRequests, "DNSMaxClientRequests", "Max DNS client requests '[data.remote]'"},
	{DNSMaxNamesResolvedIP, "DNSMaxNamesResolvedIP", "Max DNS names resolved to '[data.resolved]' by '[data.remote]'"},

	//Netfilter queue
	{NFQListedIP, "NFQListedIP", "IP [data.listed] in traffic [data.ipsrc] -> [data.ipdst]"},
	{NFQUnlistedIP, "NFQUnlistedIP", "IP [data.listed] in traffic [data.ipsrc] -> [data.ipdst]"},

	//TLS
	{TLSListedSNI, "TLSListedSNI", "SNI [data.listed] in connection [data.ipsrc] -> [data.ipdst]"},
	{TLSUnlistedSNI, "TLSUnlistedSNI", "SNI [data.listed] in connection [data.ipsrc] -> [data.ipdst]"},
	{TLSInvalidCerts, "TLSInvalidCerts", "Invalid certs validating SNI [data.sni] with subject [data.subject] in connection [data.ipsrc] -> [data.ipdst]"},
	{TLSNonTLSData, "TLSNonTLSData", "Not TLS data in connection [data.ipsrc] -> [data.ipdst]"},
}

type item struct {
	code event.Code
	name string
	desc string
}

func init() {
	for _, reg := range registry {
		event.RegisterCode(reg.code, reg.name, reg.desc)
	}
	registry = nil
}
