package main

import (
	"log"
	"bytes"
	"net/http"
	"encoding/json"
)

type Peer struct {
	url string
	uuid string
}

func http_post(url string, body interface{}) (map[string]interface{}, error) {
	data,err := json.Marshal(body)
	if err != nil {
		return nil,err
	}

	resp,err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil,err
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res,nil
}

func clusterRegisterWithPeer(addr string) { // clusterCallRegisterPeer or apiCallRegisterPeer
	p := &apiClusterPeersPOST{
		Address: web_listen,
		Token: cluster_token,
		Uuid: local_uuid,
	}

	_,err := http_post("http://" + addr + "/api/v1/cluster/peers", p)
	if err != nil {
		log.Print(err)
	}
}

func clusterRegisterPeer(addr string, uuid string, token string) {
	log.Printf("Register: %s %s %s",addr,uuid,token)
}
