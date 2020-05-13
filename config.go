package main

import (
	"log"

	"github.com/spf13/viper"
)

// config contains the complete tool configuration
type config struct {
	HTTP        httpConfig
	BackendName string
	DNS         dnsConfig
	Clients     []clientConfig

	backendViperConfig *viper.Viper
}

// httpConfig is the sub-config for the web server
type httpConfig struct {
	Listen      string
	Port        int
	BehindProxy bool
	TLS         tlsConfig
}

type tlsConfig struct {
	Cert string
	Key  string
}

type dnsConfig struct {
	Domain string
	TTL    int
}

// clientConfig is the sub-config of one single client (name + key)
type clientConfig struct {
	Name string
	Key  string
}

// getConfig reads the configuration file or exits with a fatal error
func getConfig(configfile string) *config {
	// default configuration
	cfg := config{
		HTTP: httpConfig{
			Listen:      "",
			Port:        8000,
			BehindProxy: false,
		},
		DNS: dnsConfig{
			TTL: 60,
		},
	}

	// read config file with viper
	viperConfig := viper.New()
	if configfile == "" {
		viperConfig.SetConfigName("goddns")
		viperConfig.SetConfigType("yaml")
		viperConfig.AddConfigPath("/etc/goddns/")
		viperConfig.AddConfigPath(".")
	} else {
		viperConfig.SetConfigType("yaml")
		viperConfig.SetConfigFile(configfile)
	}

	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Fatal("Cannot read configuration file: ", err)
	}

	// merge viper config into default config
	viperConfig.Unmarshal(&cfg)

	// validation (there are settings without default values, check values, ...)
	if cfg.DNS.Domain == "" {
		log.Fatal("Domain missing in configuration")
	}

	if cfg.DNS.TTL <= 0 {
		log.Fatal("Invalid DNS TTL in configuration")
	}

	if len(cfg.Clients) == 0 {
		// only a warning as we can run
		log.Print("Warning: no clients configured")
	}

	// all clients need a name and key
	for _, client := range cfg.Clients {
		if client.Name == "" {
			log.Fatal("Client without name configured")
		}

		if client.Key == "" {
			log.Fatal("Client without key configured")
		}
	}

	// finally, we need to handle the backend name and configuration
	// the backend config is not validated here
	cfg.BackendName = viperConfig.GetString("backend.name")
	if cfg.BackendName == "" {
		log.Fatal("Backend name missing in configuration")
	}
	cfg.backendViperConfig = viperConfig.Sub("backend")

	return &cfg
}
