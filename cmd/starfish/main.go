// Package main defines a program which runs a server for the starfish proxy.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bjjb/starfish"
)

func main() {
	var bindAddr, cfgFile string
	var testAndExit bool

	// Set up the command-line flags
	flag.StringVar(&bindAddr, "h", env("HTTP", ":http"), "server bind address")
	flag.StringVar(&cfgFile, "f", env("CFG", "starfish.cfg"), "configuration file")
	flag.BoolVar(&testAndExit, "t", false, "test the rules and exit")
	flag.Parse()

	if testAndExit {
		// Load the rules from the config
		_, err := loadFile(cfgFile)
		if err == nil {
			exit(0)
		}
		fmt.Fprint(stderr, err)
		exit(1)
	}

	// Otherwise we need to start a server
	if err := startServer(bindAddr, cfgFile); err != nil {
		fmt.Fprint(stderr, err)
		exit(1)
	}
}

var args = os.Args
var stdout = os.Stdout
var stderr = os.Stderr
var exit = os.Exit

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func startServer(addr, cfgFile string) error {
	handler := new(starfish.Router)
	rules, err := loadFile(cfgFile)

	handler.Replace(rules)
	handler.Push(api(""))

	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      handler,
	}

	fmt.Fprintf(stdout, "Using configuration from %s\n", cfgFile)
	fmt.Fprintf(stdout, "Starting server on %s\n", addr)
	return server.ListenAndServe()
}

func loadFile(configFile string) ([]starfish.Route, error) {
	rules, err := parseConfigFile(configFile)
	if err != nil {
		return nil, err
	}
	routes := []starfish.Route{}
	for _, rule := range rules {
		route := build(rule)
		routes = append(routes, route)
	}
	return routes, nil
}
