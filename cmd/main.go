package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"cache-proxy/internal/proxy"
)

func main() {
	port := flag.Int("port", 8080, "port on which the caching proxy server will run")
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
