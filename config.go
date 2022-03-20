package main

import (
	"flag"
)

var webListen string
var clusterListen string

func parseArgs() {
	flag.StringVar(&webListen, "web.listen", ":5838", " Address to listen on for the web interface and API")
	flag.StringVar(&clusterListen, "cluster.listen", ":5837", "Listen address for cluster, set to empty string to disable HA mode")
	flag.Parse()
}
