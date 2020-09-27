package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

var urls map[string]string
var configFile = flag.String("urls", "/etc/shorts.yml", "path to the config file with URLs (map of short: long)")
var port = flag.Int("port", 0, "port to bind to. Defaults to 0 (dynamic), so you will have to check the output to see which port was dynamically assigned.")

func readURLs(configFile string) error {
	fmt.Printf("Reading URLs from %s\n", configFile)

	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		return fmt.Errorf("Unable to read config file: #%v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &urls); err != nil {
		return fmt.Errorf("Unable to parse config file %s: #%v", configFile, err)
	}

	for short, long := range urls {
		fmt.Printf("%s => %s\n", short, long)
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	long := urls[path]

	if long == "" {
		http.Error(w, "Not found", http.StatusNotFound)

	} else {
		http.Redirect(w, r, long, http.StatusSeeOther)
	}
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

	http.HandleFunc("/", handler)

	if err := http.Serve(listener, nil); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
