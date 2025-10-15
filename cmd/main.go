package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/AMEY-GAIKAR/cache-proxy/internal/proxy"
)

const PORT = 8080

func main() {
	port := flag.Int("port", PORT, "port on which the cache-proxy server will run")
	origin := flag.String("origin", "", "url of the server to which requests will be forwarded")
	clearCache := flag.Bool("clear--cache", false, "clear the cache")
	flag.Parse()

	if *origin == "" {
		flag.Usage()
		log.Fatal("Origin URL unspecified")
	}

	proxy := proxy.InitProxy(*origin)

	log.Printf("Started server on port %d\n", *port)
	log.Printf("Forwarding requests to %s\n", *origin)

	if *clearCache {
		proxy.ClearCache()
	}

	http.Handle("/", proxy)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
