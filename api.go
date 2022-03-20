package main

import (
	"log"
	"net/http"
	"encoding/json"
)

type apiClusterPeersPOST struct {
	Address string `json:"address"`
	Token string `json:"token"`
	Uuid string `json:"uuid"`
}

func apiClusterPeers(w http.ResponseWriter, r *http.Request) {
	// args := r.URL.Query()
	log.Print(r.Method, r.RemoteAddr)
	if r.Method == "GET" {

	}
	//
	if r.Method == "POST" {
		var p apiClusterPeersPOST
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		clusterRegisterPeer(p.Address, p.Token, p.Uuid)
	}
	// target, ok := args["target"]
	// if !ok {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	//
	// metrics, ok := collectMetrics(target[0])
	// if !ok {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Fprint(w, metrics)
}
