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
		{int(xlist.MD5), true},
		{int(xlist.SHA1), true},
		{int(xlist.SHA256), true},
		//invalid values as the time of the writting if the test:
		{-1, false},
		{6, false},
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
		{xlist.MD5, "md5"},
		{xlist.SHA1, "sha1"},
		{xlist.SHA256, "sha256"},
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
		//test md5
		{"c9ef1a05668261b882e9267af006f78d", xlist.MD5, true},
		{"C9EF1A05668261B882e9267af006f78d", xlist.MD5, true},
		{"nidecona", xlist.MD5, false},
		{"x9ef1a05668261b882e9267af006f78d", xlist.MD5, false},
		{"0c9ef1a05668261b882e9267af006f78d", xlist.MD5, false},
		//test sha1
		{"4544f891cb3c190366bc5d0d331ae17e254b26e6", xlist.SHA1, true},
		{"4544f891CB3C190366BC5d0D331AE17E254B26E6", xlist.SHA1, true},
		{"nidecona", xlist.SHA1, false},
		{"X544f891CB3C190366BC5d0D331AE17E254B26E6", xlist.SHA1, false},
		//test sha256
		{"00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, true},
		{"00015b14C28C2951f6D628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, true},
		{"nidecona", xlist.SHA256, false},
		{"X0015b14C28C2951f6D628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, false},

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
		//test md5
		{"c9ef1a05668261b882e9267af006f78d", xlist.MD5},
		//test sha1
		{"4544f891cb3c190366bc5d0d331ae17e254b26e6", xlist.SHA1},
		//test sha256
		{"00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256},
		//unexpected
		{"-12.34.23.", xlist.Resource(-1)},
	}
	corder := []xlist.Resource{
		xlist.IPv4, xlist.IPv6,
		xlist.MD5, xlist.SHA1, xlist.SHA256,
		xlist.Domain,
	}
	for _, test := range tests {
		got, err := xlist.ResourceType(test.name, corder)
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
		{"c9ef1a05668261b882e9267af006f78d", xlist.MD5, nil},
		{"4544f891cb3c190366bc5d0d331ae17e254b26e6", xlist.SHA1, nil},
		{"00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, nil},
		// not valid
		{"12.11", xlist.IPv4, xlist.ErrBadRequest},
		{"12.11", xlist.IPv6, xlist.ErrBadRequest},
		{"-www.com", xlist.Domain, xlist.ErrBadRequest},
		{"X9ef1a05668261b882e9267af006f78d", xlist.MD5, xlist.ErrBadRequest},
		{"X544f891cb3c190366bc5d0d331ae17e254b26e6", xlist.SHA1, xlist.ErrBadRequest},
		{"X0015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, xlist.ErrBadRequest},

		//unexpected
		{"12.34.23.2", xlist.Resource(-1), xlist.ErrBadRequest},
		{"12.34.23.3", xlist.Resource(10), xlist.ErrBadRequest},
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
		{"c9ef1a05668261b882e9267af006f78d", xlist.MD5, true, "c9ef1a05668261b882e9267af006f78d"},
		{"4544f891cb3c190366bc5d0d331ae17e254b26e6", xlist.SHA1, true, "4544f891cb3c190366bc5d0d331ae17e254b26e6"},
		{"00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, true, "00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8"},
		// not canonical
		{"fe80::3289:ad8e:8259:c878", xlist.IPv6, true, "fe80::3289:ad8e:8259:c878"},
		{"FE80::3289:AD8E:8259:c878", xlist.IPv6, true, "fe80::3289:ad8e:8259:c878"},
		{"fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:0050", xlist.IPv6, true, "fd8c:15c7:33f2:ed00:b5cb:bbdf:8266:50"},
		{"fd8c:15c7:33f2:ed00::bbdf:8266:fa50", xlist.IPv6, true, "fd8c:15c7:33f2:ed00:0:bbdf:8266:fa50"},
		{"fd8c:15c7:33f2:0000:0000:bbdf:8266:fa50", xlist.IPv6, true, "fd8c:15c7:33f2::bbdf:8266:fa50"},
		{"WWW.DOMINIO.com", xlist.Domain, true, "www.dominio.com"},
		{"C9EF1A05668261B882E9267AF006F78D", xlist.MD5, true, "c9ef1a05668261b882e9267af006f78d"},
		{"4544f891CB3C190366BC5D0d331ae17e254b26e6", xlist.SHA1, true, "4544f891cb3c190366bc5d0d331ae17e254b26e6"},
		{"00015B14C28C2951F6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8", xlist.SHA256, true, "00015b14c28c2951f6d628098ce6853e14300f1b7d6d985e18d508f9807f44d8"},

		// not valid
		{"12.11", xlist.IPv4, false, "12.11"},
		{"12.11", xlist.IPv6, false, "12.11"},
		{"-www.com", xlist.Domain, false, "-www.com"},
		{"X9ef1a05668261b882e9267af006f78d", xlist.MD5, false, "X9ef1a05668261b882e9267af006f78d"},
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
	if err != xlist.ErrBadRequest {
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
			[]xlist.Resource{xlist.IPv4, xlist.IPv6}},
		{[]xlist.Resource{xlist.IPv6, xlist.IPv4, xlist.IPv6, xlist.Domain, xlist.Domain},
			[]xlist.Resource{xlist.IPv4, xlist.IPv6, xlist.Domain}},
		{[]xlist.Resource{xlist.IPv4, xlist.MD5, xlist.IPv6, xlist.MD5, xlist.IPv4},
			[]xlist.Resource{xlist.IPv4, xlist.IPv6, xlist.MD5}},
		{[]xlist.Resource{xlist.IPv4, xlist.SHA256, xlist.IPv6, xlist.SHA1, xlist.IPv4, xlist.SHA256},
			[]xlist.Resource{xlist.IPv4, xlist.IPv6, xlist.SHA1, xlist.SHA256}},
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
