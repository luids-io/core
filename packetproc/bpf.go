// Copyright 2018 Luis Guillén Civera <luisguillenc@gmail.com>. All rights reserved.

package packetproc

import (
	"fmt"
	"net"
	"strings"
)

// BPFilter is a simple helper for constructing bpf filters
type BPFilter struct {
	//Base is base filter
	Base string
	//Include is a list of CIDR or IP to include in the filter
	Include []string
	//Exclude is a list of CIDR or IP to exclude from the filter
	Exclude []string
}

// Filter returns the filter in bpf format
func (bpf BPFilter) Filter() (string, error) {
	filter := bpf.Base
	include, err := includeFilter(bpf.Include)
	if err != nil {
		return "", err
	}
	if include != "" {
		filter = filter + " and " + include
	}
	exclude, err := excludeFilter(bpf.Exclude)
	if err != nil {
		return "", err
	}
	if exclude != "" {
		filter = filter + " and " + exclude
	}
	return filter, nil
}

func parseIPsNets(items []string) ([]string, []string, error) {
	ips := make([]string, 0)
	nets := make([]string, 0)
	for _, item := range items {
		_, _, err := net.ParseCIDR(item)
		if err == nil {
			nets = append(nets, item)
		} else {
			ip := net.ParseIP(item)
			if ip != nil {
				ips = append(ips, item)
			} else {
				return nil, nil, fmt.Errorf("%s not ip or CIDR", item)
			}
		}
	}
	return ips, nets, nil
}

func includeFilter(exclude []string) (string, error) {
	ips, nets, err := parseIPsNets(exclude)
	if err != nil {
		return "", err
	}
	filter := ""
	if len(ips) > 0 {
		s := fmt.Sprintf("host (%s)", strings.Join(ips, " or "))
		filter = filter + s
	}
	if len(nets) > 0 {
		if filter != "" {
			filter = filter + " or "
		}
		s := fmt.Sprintf("net (%s)", strings.Join(nets, " or "))
		filter = filter + s
	}
	if filter != "" {
		filter = fmt.Sprintf("(%s)", filter)
	}
	return filter, nil
}

func excludeFilter(exclude []string) (string, error) {
	ips, nets, err := parseIPsNets(exclude)
	if err != nil {
		return "", err
	}
	filter := ""
	if len(ips) > 0 {
		s := fmt.Sprintf("not host (%s)", strings.Join(ips, " or "))
		filter = filter + s
	}
	if len(nets) > 0 {
		if filter != "" {
			filter = filter + " and "
		}
		s := fmt.Sprintf("not net (%s)", strings.Join(nets, " or "))
		filter = filter + s
	}
	if filter != "" {
		filter = fmt.Sprintf("(%s)", filter)
	}
	return filter, nil
}
