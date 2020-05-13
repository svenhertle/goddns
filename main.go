package main

import (
	"flag"
	"fmt"
)

var goddnsVersion = "0.1.0"

func main() {
	// command line flags
	var configfile = flag.String("config", "", "load configuration from this file")
	var showVersion = flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Println(goddnsVersion)
	} else {
		NewGoDDNS(*configfile)
	}
}
