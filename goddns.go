package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/svenhertle/goddns/backends"
)

// GoDDNS is a dynamic DNS server/updater
type GoDDNS struct {
	cfg *config

	backend backends.Backend
}

// NewGoDDNS creates a new GoDDNS instance
func NewGoDDNS(configfile string) *GoDDNS {
	goddns := GoDDNS{}

	// read config
	goddns.cfg = getConfig(configfile)

	// initialize backend
	goddns.backend = backends.GetBackend(goddns.cfg.BackendName, goddns.cfg.backendViperConfig,
		goddns.cfg.DNS.Domain, goddns.cfg.DNS.TTL)

	// start web server
	goddns.startHTTP(goddns.cfg.HTTP.Listen, goddns.cfg.HTTP.Port)

	return &goddns
}

// startHTTP configured and starts the web server
func (goddns *GoDDNS) startHTTP(listen string, port int) {
	// routes
	http.HandleFunc("/", goddns.handleBase)
	http.HandleFunc("/v1/", goddns.handleVersion1)

	// listening address and port
	listenAddress := fmt.Sprintf("[%s]:%d", listen, port)
	log.Println("Listening on", listenAddress)

	// run
	var err error
	if goddns.cfg.HTTP.TLS.Cert != "" && goddns.cfg.HTTP.TLS.Key != "" {
		// with TLS
		err = http.ListenAndServeTLS(listenAddress, goddns.cfg.HTTP.TLS.Cert, goddns.cfg.HTTP.TLS.Key, nil)
	} else {
		// without TLS
		err = http.ListenAndServe(listenAddress, nil)
	}
	if err != nil {
		log.Fatal("Starting the web server failed: ", err)
	}
}

// handleBase handles HTTP requests on "/"
func (goddns *GoDDNS) handleBase(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

// handleVersion1 handles HTTP requests by clients (v1)
func (goddns *GoDDNS) handleVersion1(w http.ResponseWriter, r *http.Request) {
	// get key
	key := r.URL.Query().Get("key")
	ip := r.URL.Query().Get("ip")

	// check if key is given
	if key == "" {
		http.Error(w, "key missing", http.StatusBadRequest)
		return
	}

	// check if key is valid for a client
	var name string
	for _, client := range goddns.cfg.Clients {
		if client.Key == key {
			name = client.Name
			break
		}
	}

	if name == "" {
		http.Error(w, "key unknown", http.StatusBadRequest)
		return
	}

	// get remote ip if we did not get it
	if ip == "" {
		if goddns.cfg.HTTP.BehindProxy {
			ip = GetIPBehindProxy(r)
		} else {
			ip = GetIPDirect(r)
		}

		if ip == "" {
			http.Error(w, "internal error", http.StatusInternalServerError)
			log.Println("Cannot determine IP address of client")
			return
		}
	}

	// validate ip
	addressType, ipValid := CheckIP(ip)
	if !ipValid {
		http.Error(w, "ip invalid", http.StatusBadRequest)
		log.Printf("Got invalid IP address %s for client %s\n", ip, name)
		return
	}

	// run the update
	log.Printf("New IP for %s: %s\n", name, ip)
	err := goddns.backend.Update(name, ip, addressType)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Println("Error while updating an entry:", err)
		return
	}

	fmt.Fprintf(w, "ok\n")
}
