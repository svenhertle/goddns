package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/svenhertle/goddns/backends"
)

// GetIPDirect reads the client IP directly from the HTTP request.
func GetIPDirect(r *http.Request) string {
	tmp := r.RemoteAddr // format: ip:port

	// remove port
	tmp = tmp[0:strings.LastIndex(tmp, ":")]

	// remove "[" and "]" (for ipv6)
	tmp = strings.Replace(tmp, "[", "", 1)
	tmp = strings.Replace(tmp, "]", "", 1)

	return tmp
}

// GetIPBehindProxy checks typical HTTP headers set by proxies for the client IP address.
func GetIPBehindProxy(r *http.Request) string {
	// check the following headers
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, header := range headers {
		value := r.Header.Get(header)
		if value != "" {
			// found, we ware done
			return value
		}
	}

	return ""
}

// CheckIP validates the IP address and returns its `AddressType` and whether it is valid or not.
func CheckIP(ip string) (backends.AddressType, bool) {
	// parse IP
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "", false
	}

	// get address type
	if parsed.To4() != nil {
		return backends.IPv4, true
	}
	return backends.IPv6, true
}
