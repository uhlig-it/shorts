package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"go.yaml.in/yaml/v4"
)

var urls map[string]string
var configFile = flag.String("urls", "/etc/shorts.yml", "path to the config file with URLs (map of short: long)")
var port = flag.Int("port", 0, "port to bind to. Defaults to 0 (dynamic), so you will have to check the output to see which port was dynamically assigned.")

func readURLs(configFile string) error {
	fmt.Printf("Reading URLs from %s\n", configFile)

	// #nosec G304 -- file path comes from the -urls flag; reading arbitrary paths is intended.
	yamlFile, err := os.ReadFile(configFile)

	if err != nil {
		return fmt.Errorf("unable to read config file: #%v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &urls); err != nil {
		return fmt.Errorf("unable to parse config file %s: #%v", configFile, err)
	}

	for short, long := range urls {
		fmt.Printf("%s => %s\n", short, long)
	}

	return nil
}

func main() {
	flag.Parse()
	err := readURLs(*configFile)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	bindAddress := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", bindAddress)

	if err != nil {
		fmt.Printf("Unable to start listening on %s #%v\n", bindAddress, err)
		os.Exit(1)
	}

	fmt.Printf("Listening at %s\n", listener.Addr())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		short := r.URL.Path[1:]
		long := urls[short]

		if long == "" {
			http.Error(w, "Not found", http.StatusNotFound)

		} else {
			http.Redirect(w, r, long, http.StatusSeeOther)
		}
	})

	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
