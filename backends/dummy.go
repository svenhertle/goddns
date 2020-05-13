package backends

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// DummyBackend is a GoDDNS backend for debugging, it just prints the DNS record changes.
type DummyBackend struct {
	domain string
	ttl    int
}

// Configure does not have much to do for the Dummy backend
func (backend *DummyBackend) Configure(cfg *viper.Viper, domain string, ttl int) error {
	// store dns settings
	backend.domain = domain
	backend.ttl = ttl

	return nil
}

// Update writes a new record to the terminal
func (backend *DummyBackend) Update(name string, ip string, addressType AddressType) error {
	fqdn := fmt.Sprintf("%s.%s", name, backend.domain)
	log.Printf("Change record: %s -> %s (Type: %s, TTL: %d)\n", fqdn, ip, addressType, backend.ttl)

	return nil
}
