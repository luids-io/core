// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package xlist_test

import (
	"context"
	"testing"

	"github.com/luids-io/core/xlist"
)

func TestResourceIsValid(t *testing.T) {
	var tests = []struct {
		input int
		want  bool
	}{
		{int(xlist.IPv4), true},
		{int(xlist.IPv6), true},
		{int(xlist.Domain), true},
		//invalid values as the time of the writting if the test:
		{-1, false},
		{3, false},
		{10, false},
	}
	for _, test := range tests {
		resource := xlist.Resource(test.input)
		if got := resource.IsValid(); got != test.want {
			t.Errorf("resource[%v].IsValid() = %v", test.input, got)
		}
	}
}

func TestResourceString(t *testing.T) {
	var tests = []struct {
		input xlist.Resource
		want  string
	}{
		{xlist.IPv4, "ip4"},
		{xlist.IPv6, "ip6"},
		{xlist.Domain, "domain"},
		{xlist.Resource(-1), "unkown(-1)"},
	}
	for _, test := range tests {
		if got := test.input.String(); got != test.want {
			t.Errorf("resource[%v].String() = %v", test.input, got)
		}
	}
}

func TestResourceInArray(t *testing.T) {
	var tests = []struct {
		resource xlist.Resource
		array    []xlist.Resource
		want     bool
	}{
		{xlist.IPv4, []xlist.Resource{}, false},
		{xlist.IPv4, []xlist.Resource{xlist.Domain, xlist.IPv6}, false},
		{xlist.IPv4, []xlist.Resource{xlist.Domain, xlist.IPv6, xlist.IPv4}, true},
		{xlist.IPv4, []xlist.Resource{xlist.IPv4, xlist.IPv4, xlist.Domain}, true},
		//invalid values as the time of the writting if the test:
		{xlist.IPv4, []xlist.Resource{xlist.Resource(-1)}, false},
		{xlist.IPv4, []xlist.Resource{xlist.Resource(-1), xlist.IPv4}, true},
		{xlist.Resource(-1), []xlist.Resource{xlist.Resource(-1), xlist.IPv4}, true},
	}
	for _, test := range tests {
		resource := xlist.Resource(test.resource)
		if got := resource.InArray(test.array); got != test.want {
			t.Errorf("resource[%v].InArray(%v) = %v", test.resource, test.array, got)
		}
	}
}

func TestValidResource(t *testing.T) {
	var tests = []struct {
		name     string
		resource xlist.Resource
		want     bool
	}{
		//test ipv4
		{"12.34.23.1", xlist.IPv4, true},
		{"12.34.23.256", xlist.IPv4, false},
		{"12.34.23.", xlist.IPv4, false},
		{"fe80::3289:ad8e:8259:c878", xlist.IPv4, false},
		{"nombre.com", xlist.IPv4, false},
		//test ip6
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6, true},
		{"fe80:3289:ad8e:8259:c878", xlist.IPv6, false},
		//test domain
		{"dominio", xlist.Domain, true},
		{"dominio.com", xlist.Domain, true},
		{"www.dominio.com", xlist.Domain, true},
		{"-sdf.com", xlist.Domain, false},
		//unexpected
		{"kk.com", xlist.Resource(-1), false},
	}
	for _, test := range tests {
		if got := xlist.ValidResource(test.name, test.resource); got != test.want {
			t.Errorf("ValidResource(%v, %v) = %v", test.name, test.resource, got)
		}
	}
}

func TestResourceType(t *testing.T) {
	var tests = []struct {
		name string
		want xlist.Resource
	}{
		//TODO: add more testing
		//test ipv4
		{"12.34.23.1", xlist.IPv4},
		//test ip6
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6},
		//test domain
		{"www.dominio.com", xlist.Domain},
		//unexpected
		{"-12.34.23.", xlist.Resource(-1)},
	}
	for _, test := range tests {
		got, err := xlist.ResourceType(test.name, xlist.Resources)
		if test.want == xlist.Resource(-1) {
			if err == nil {
				t.Errorf("ResourceType(%v) = %v, %v", test.name, got, err)
			}
		} else {
			if got != test.want {
				t.Errorf("ResourceType(%v) = %v, %v", test.name, got, err)
			}
		}
	}
}

func TestDoValidation(t *testing.T) {
	var tests = []struct {
		name     string
		resource xlist.Resource
		want     error
	}{
		{"12.34.23.1", xlist.IPv4, nil},
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6, nil},
		{"www.dominio.com", xlist.Domain, nil},
		// not valid
		{"12.11", xlist.IPv4, xlist.ErrBadResourceFormat},
		{"12.11", xlist.IPv6, xlist.ErrBadResourceFormat},
		{"-www.com", xlist.Domain, xlist.ErrBadResourceFormat},
		//unexpected
		{"12.34.23.2", xlist.Resource(-1), xlist.ErrResourceNotSupported},
		{"12.34.23.3", xlist.Resource(10), xlist.ErrResourceNotSupported},
	}
	for _, test := range tests {
		_, _, got := xlist.DoValidation(context.Background(), test.name, test.resource, true)
		if got != test.want {
			t.Errorf("DoValidation(ctx, %v, %v, true) = err(%v)", test.name, test.resource, got)
		}
	}
}

func TestCanonicalize(t *testing.T) {
	var tests = []struct {
		name     string
		resource xlist.Resource
		want     bool
		wantName string
	}{
		{"12.34.23.1", xlist.IPv4, true, "12.34.23.1"},
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6, true, "fe80::3289:ad8e:8259:c878"},
		{"fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:fa50", xlist.IPv6, true, "fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:fa50"},
		{"www.dominio.com", xlist.Domain, true, "www.dominio.com"},
		// not canonical
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6, true, "fe80::3289:ad8e:8259:c878"},
		{"FE80::3289:AD8E:8259:c878", xlist.IPv6, true, "fe80::3289:ad8e:8259:c878"},
		{"fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:0050", xlist.IPv6, true, "fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:50"},
		{"fd8c:15c7:33f2:ed00::bbdf:8266:fa50", xlist.IPv6, true, "fd8c:15c7:33f2:ed00:0:bbdf:8266:fa50"},
		{"fd8c:15c7:33f2:0000:0000:bbdf:8266:fa50", xlist.IPv6, true, "fd8c:15c7:33f2::bbdf:8266:fa50"},
		{"WWW.DOMINIO.com", xlist.Domain, true, "www.dominio.com"},
		// not valid
		{"12.11", xlist.IPv4, false, "12.11"},
		{"12.11", xlist.IPv6, false, "12.11"},
		{"-www.com", xlist.Domain, false, "-www.com"},
	}
	for _, test := range tests {
		gotName, got := xlist.Canonicalize(test.name, test.resource)
		if got != test.want || gotName != test.wantName {
			t.Errorf("Canonicalize(%v, %v) = (%v, %v)", test.name, test.resource, got, gotName)
		}
	}
}

func TestDoValidationContext(t *testing.T) {
	var err error
	ctx := context.Background()
	name := "127.0.0.1"
	resource := xlist.IPv4
	// first validation
	_, ctx, err = xlist.DoValidation(ctx, name, resource, false)
	if err != nil {
		t.Errorf("DoValidation(background, %v, %v, false) = err(%v)", name, resource, err)
	}
	// new validations in the context must return without check...
	_, ctx, err = xlist.DoValidation(ctx, name, resource, false)
	if err != nil {
		t.Errorf("DoValidation(background, %v, %v, false) = err(%v)", name, resource, err)
	}
	// we set an invalid value for the resource type...
	name = "www.host.com"
	if err != nil {
		t.Errorf("DoValidation(background, %v, %v, false) = err(%v)", name, resource, err)
	}
	// now we test with force enabled...
	_, ctx, err = xlist.DoValidation(ctx, name, resource, true)
	if err != xlist.ErrBadResourceFormat {
		t.Errorf("DoValidation(ctx, %v, %v, true) = err(%v)", name, resource, err)
	}
}

func TestClearResourceDups(t *testing.T) {
	var tests = []struct {
		in  []xlist.Resource
		out []xlist.Resource
	}{
		{[]xlist.Resource{}, []xlist.Resource{}},
		{[]xlist.Resource{xlist.Resource(-1)}, []xlist.Resource{}},
		{[]xlist.Resource{xlist.IPv4}, []xlist.Resource{xlist.IPv4}},
		{[]xlist.Resource{xlist.IPv4, xlist.IPv4, xlist.Resource(-1)},
			[]xlist.Resource{xlist.IPv4}},
		{[]xlist.Resource{xlist.IPv4, xlist.IPv6, xlist.IPv4},
			[]xlist.Resource{xlist.IPv4, xlist.IPv6}},
		{[]xlist.Resource{xlist.IPv6, xlist.IPv4, xlist.IPv6},
			[]xlist.Resource{xlist.IPv6, xlist.IPv4}},
		{[]xlist.Resource{xlist.IPv6, xlist.IPv4, xlist.IPv6, xlist.Domain, xlist.Domain},
			[]xlist.Resource{xlist.IPv6, xlist.IPv4, xlist.Domain}},
	}
	for _, test := range tests {
		got := xlist.ClearResourceDups(test.in)
		if !cmpResourceSlice(test.out, got) {
			t.Errorf("ClearResourceDups(%v) = %v", test.in, got)
		}
	}
}

func cmpResourceSlice(a, b []xlist.Resource) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
