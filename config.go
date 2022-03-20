package main

import (
	"log"
	"flag"
	"github.com/google/uuid"
)

var web_listen string
var cluster_listen string
var cluster_token string
var local_uuid string

func parseArgs() {
	flag.StringVar(&web_listen, "web.listen", ":5838", " Address to listen on for the web interface and API")
	flag.StringVar(&cluster_listen, "cluster.listen", ":5837", "Listen address for cluster, set to empty string to disable HA mode")
	flag.Parse()

	local_uuid = uuid.New().String()
	cluster_token = "mytoken"

	log.Printf("configuration: web.listen = %s", web_listen);
	log.Printf("configuration: cluster.listen = %s", cluster_listen);
	log.Printf("configuration: local.uuid = %s", local_uuid);
}
