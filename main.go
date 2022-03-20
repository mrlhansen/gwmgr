package main

import (
	"log"
	"net/http"
)

func main() {
	parseArgs()
	go discoverStart()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/api/v1/cluster/peers", apiClusterPeers)
	err := http.ListenAndServe(web_listen, nil)
	log.Fatal(err)
}
