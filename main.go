package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"gopkg.in/yaml.v2"
)

var urls map[string]string
var configFile = flag.String("urls", "/etc/shorts.yml", "path to the config file with URLs (map of short: long)")
var port = flag.Int("port", 0, "port to bind to. Defaults to 0 (dynamic), so you will have to check the output to see which port was dynamically assigned.")

func readURLs(configFile string) error {
	log.Printf("Reading URLs from %s\n", configFile)

	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		return fmt.Errorf("Unable to read config file: #%v", err)
	}

	if err := yaml.Unmarshal(yamlFile, &urls); err != nil {
		return fmt.Errorf("Unable to parse config file %s: #%v", configFile, err)
	}

	for short, long := range urls {
		log.Printf("%s => %s", short, long)
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
		log.Fatalf(err.Error())
	}

	bindAddress := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", bindAddress)

	if err != nil {
		log.Fatalf("Unable to start listening on %s #%v", bindAddress, err)
	}

	log.Printf("Listening at %s", listener.Addr())

	http.HandleFunc("/", handler)
	log.Fatal(http.Serve(listener, nil))
}
