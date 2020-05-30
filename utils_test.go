package main

import (
	"net/http"
	"testing"

	"github.com/svenhertle/goddns/backends"
)

func TestGetIPDirect(t *testing.T) {
	var testdata = []struct {
		in  string
		out string
	}{
		{"10.0.0.1:1234", "10.0.0.1"},
		{"[2000::1]:1234", "2000::1"},
		{"[10.0.0.1]:1234", "10.0.0.1"},
	}

	for _, d := range testdata {
		r := http.Request{
			RemoteAddr: d.in,
		}

		got := GetIPDirect(&r)
		if got != d.out {
			t.Errorf("GetIPDirect: input %s, got %s, want %s", d.in, got, d.out)
		}
	}
}

func TestGetIPBehindProxy(t *testing.T) {
	var testdata = []struct {
		header string
		value  string
		out    string
	}{
		{"X-Forwarded-For", "10.0.0.1", "10.0.0.1"},
		{"X-Real-IP", "10.0.0.1", "10.0.0.1"},
		{"X-Forwarded-For", "2000::1", "2000::1"},
		{"X-Real-IP", "2000::1", "2000::1"},
		{"No-Real-Header", "2000::1", ""},
	}

	for _, d := range testdata {
		r := http.Request{
			Header: http.Header{},
		}
		r.Header.Set(d.header, d.value)

		got := GetIPBehindProxy(&r)
		if got != d.out {
			t.Errorf("GetIPBehindProxy: input %s: %s, got %s, want %s", d.header, d.value, got, d.out)
		}
	}
}

func TestCheckIP(t *testing.T) {
	var testdata = []struct {
		in             string
		outAddressType backends.AddressType
		outValid       bool
	}{
		{"10.0.0.1", backends.IPv4, true},
		{"10.0.0", "", false},
		{"10.0.0.1.2", "", false},

		{"2000::1", backends.IPv6, true},
		{"::1", backends.IPv6, true},
		{"2000::1::2", "", false},
		{":1", "", false},

		{"", "", false},
	}

	for _, d := range testdata {
		addressType, valid := CheckIP(d.in)
		if addressType != d.outAddressType || valid != d.outValid {
			t.Errorf("GetIPDirect: input %s; got %s,%t; want %s,%t", d.in, addressType, valid, d.outAddressType, d.outValid)
		}
	}
}
