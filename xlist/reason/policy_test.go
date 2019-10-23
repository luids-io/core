// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package reason_test

import (
	"testing"

	"github.com/luids-io/core/xlist/reason"
)

func TestReasonClean(t *testing.T) {
	var tests = []struct {
		in   string
		want string
	}{
		{"[policy][/policy]razon", "razon"},
		{"[policy][/policy][/policy]razon", "[/policy]razon"},
		{"[policy][/policy][policy]razon", "[policy]razon"},
		{"ra[policy][/policy]zon", "razon"},
		{"ra[POLICY][/policy]zon", "razon"},
		{"[POLICY][/policy][policy][/policy]razon", "razon"},
	}
	for _, test := range tests {
		if got := reason.Clean(test.in); got != test.want {
			t.Errorf("Clean(%v) = %v", test.in, got)
		}
	}
}

func TestFromString(t *testing.T) {
	var tests = []struct {
		inStr   string
		inField string
		want    string
	}{
		{"[policy]dns=nxdomain[/policy]", "dns", "nxdomain"},
		{"[policy]dns=nxdomain,log[/policy]", "log", ""},
		{"[policy]dns= nxdomain, log =[/policy]", "dns", "nxdomain"},
		{"[policy]dns= nxdomain, log =,otro=kk[/policy]", "log", ""},
		{"[policy]dns= nxdomain, log =,otro=kk[/policy]", "otro", "kk"},
		{"[policy]dns=ip4:127.0.0.1,event=high[/policy]", "dns", "ip4:127.0.0.1"},
	}
	for _, test := range tests {
		p := reason.NewPolicy()
		err := p.FromString(test.inStr)
		if err != nil {
			t.Fatalf("FromString(%s): %v", test.inStr, err)
		}
		got, ok := p.Get(test.inField)
		if !ok {
			t.Fatalf("field %s don't loaded", test.inField)
		}
		if got != test.want {
			t.Errorf("want='%s', got='%s'", test.want, got)
		}
	}
}

func TestFromStringErr(t *testing.T) {
	var tests = []struct {
		in      string
		wantErr bool
	}{
		{"[policy]dns=nxdomain[/policy]", false},
		{"[policy]dns=nxdomain", true},
		{"[policy]dns=nxdomain[/POLICY]", false},
		{"[policy]dns-=nxdomain[/policy]", true},
		{"[policy]dns=nxdom ain[/policy]", true},
		{"[policy]dns=nxdomain[/policy]", false},
	}
	for _, test := range tests {
		p := reason.NewPolicy()
		err := p.FromString(test.in)
		if err != nil && !test.wantErr {
			t.Errorf("FromString(%s) unexpected err: %v", test.in, err)
		} else if err == nil && test.wantErr {
			t.Errorf("FromString(%s) expected err", test.in)
		}
	}
}

func TestExtractPolicy(t *testing.T) {
	var tests = []struct {
		in   string
		key  string
		want string
	}{
		{"[policy]dns=nxdomain[/policy]", "dns", "nxdomain"},
		{"[policy]dns=nxdomain,event=info[/policy]", "event", "info"},
		{"[policy]field1=aa[/policy][policy]field2=aa[/policy]", "field2", "aa"},
		{"[policy]field1=aa[/policy][policy]field2=bb[/policy]", "field2", "bb"},
		{"[policy]field1=aa, field2=aa[/policy][policy]field2=cc[/policy]", "field1", "aa"},
		{"[policy]field1=aa, field2=aa[/policy][policy]field2=cc[/policy]", "field2", "cc"},
	}
	for _, test := range tests {
		p, _, err := reason.ExtractPolicy(test.in)
		if err != nil {
			t.Errorf("ExtractPolicy(%s) unexpected err: %v", test.in, err)
		}
		if got, _ := p.Get(test.key); got != test.want {
			t.Errorf("ExtractPolicy(%s) want=%v got=%v", test.in, test.want, got)
		}
	}
}
