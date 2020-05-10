package ipfilter_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/luids-io/core/ipfilter"
)

func TestEmpty(t *testing.T) {
	var tests = []struct {
		inputLists []string
		want       bool
	}{
		{[]string{}, true},
		{[]string{"192.168.1.1"}, false},
		{[]string{"192.168.1.0/24", "127.0.0.0/8"}, false},
		{[]string{"192.168.1.X/24"}, true},
	}
	for _, test := range tests {
		filter := ipfilter.Whitelist(test.inputLists)
		if got := filter.Empty(); got != test.want {
			t.Errorf("filter.Empty() = %v", got)
		}
		filter = ipfilter.Blacklist(test.inputLists)
		if got := filter.Empty(); got != test.want {
			t.Errorf("filter.Empty() = %v", got)
		}
	}
}

func TestWhitlisting(t *testing.T) {
	var tests = []struct {
		inputLists []string
		check      string
		want       ipfilter.Action
	}{
		{[]string{}, "192.168.1.2", ipfilter.Deny},
		{[]string{"192.168.1.1"}, "192.168.1.2", ipfilter.Deny},
		{[]string{"192.168.1.2"}, "192.168.1.2", ipfilter.Accept},
		{[]string{"192.168.1.0/24", "23.3.4.3", "127.0.0.0/8"},
			"192.168.1.2", ipfilter.Accept},
		{[]string{"192.168.1.0/24", "23.3.4.3", "127.0.0.0/8"},
			"127.0.0.2", ipfilter.Accept},
		{[]string{"192.168.1.X/24", "23.3.4.3", "127.0.0.0/8"},
			"192.168.1.2", ipfilter.Deny},
	}
	for _, test := range tests {
		filter := ipfilter.Whitelist(test.inputLists)
		if got := filter.Check(net.ParseIP(test.check)); got != test.want {
			t.Errorf("Check(%v) = %v", test.check, got)
		}
	}
}

func TestBlacklisting(t *testing.T) {
	var tests = []struct {
		inputLists []string
		check      string
		want       ipfilter.Action
	}{
		{[]string{}, "192.168.1.2", ipfilter.Accept},
		{[]string{"192.168.1.1"}, "192.168.1.2", ipfilter.Accept},
		{[]string{"192.168.1.2"}, "192.168.1.2", ipfilter.Deny},
		{[]string{"192.168.1.0/24", "23.3.4.3", "127.0.0.0/8"},
			"192.168.1.2", ipfilter.Deny},
		{[]string{"192.168.1.0/24", "23.3.4.3", "127.0.0.0/8"},
			"127.0.0.2", ipfilter.Deny},
		{[]string{"192.168.1.X/24", "23.3.4.3", "127.0.0.0/8"},
			"192.168.1.2", ipfilter.Accept},
	}
	for _, test := range tests {
		filter := ipfilter.Blacklist(test.inputLists)
		if got := filter.Check(net.ParseIP(test.check)); got != test.want {
			t.Errorf("Check(%v) = %v", test.check, got)
		}
	}
}

func ExampleWhitelist() {
	list := []string{
		"192.168.1.0/24",
		"23.3.4.3",
		"127.0.0.0/8",
	}
	filter := ipfilter.Whitelist(list)

	ok := filter.Check(net.ParseIP("192.168.1.2"))
	if ok {
		fmt.Println("IP aceptada")
	}
}
