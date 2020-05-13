package backends

import (
	"log"

	"github.com/spf13/viper"
)

const (
	// IPv4 means that the DNS records is a IPv4 address
	IPv4 = "A"

	// IPv6 means that the DNS records is a IPv4 address
	IPv6 = "AAAA"
)

// AddressType is the address type of a DNS record: `IPv4` or `IPv6`.
type AddressType string

// Backend is the interface for GoDDNS backend and allows a simple extension of GoDDNS.
type Backend interface {
	Configure(cfg *viper.Viper, domain string, ttl int) error

	Update(name string, ip string, addressType AddressType) error
}

// GetBackend returns the specified backend ready for use or exits with a fatal error.
func GetBackend(name string, cfg *viper.Viper, domain string, ttl int) Backend {
	var backend Backend
	// find backend (add new here)
	if name == "powerdns" {
		backend = &PowerDNSBackend{}
	} else if name == "dummy" {
		backend = &DummyBackend{}
	} else {
		log.Fatal("Unkown backend: ", name)
	}

	// configure it
	err := backend.Configure(cfg, domain, ttl)
	if err != nil {
		log.Fatal("Failed to configure backend: ", err)
	}

	// and return it
	return backend
}
