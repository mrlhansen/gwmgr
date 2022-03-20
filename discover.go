package main

import (
	"log"
	"fmt"
	"net"
	"time"
	"regexp"
	"strings"
	"strconv"
	"encoding/binary"
)

func writeDiscoverMessage() []byte {
	msg := fmt.Sprintf("v1|%s|%s|%s", web_listen, cluster_token, local_uuid)
	return []byte(msg)
}

func readDiscoverMessage(addr string, msg []byte) {
	var web string
	var token string
	var uuid string

	list := strings.Split(string(msg), "|")
	if len(list) != 4 {
		return
	}

	web = list[1]
	token = list[2]
	uuid = list[3]

	if token != cluster_token {
		return
	}

	if uuid == local_uuid {
		// return
	}

	host,port,_ := net.SplitHostPort(web)
	if host == "" || host == "0.0.0.0" {
		host,_,_ = net.SplitHostPort(addr)
	}

	clusterRegisterWithPeer(host + ":" + port)
}

func calcBroadcast(ip string) string {
	var brd = "255.255.255.255"
	var addr uint32
	var mask uint32

	if ip == "" || ip == "0.0.0.0" {
		return brd
	}

	// Find the listen address on the host and calculate the associated
	// broadcast address from the CIDR prefix
	list,err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _,v := range list {
		s := v.String()
		if strings.HasPrefix(s, ip) {
			re := regexp.MustCompile(`(\d+).(\d+).(\d+).(\d+)/(\d+)`)
			match := re.FindStringSubmatch(s)
			if len(match) == 0 {
				log.Fatal("failed to parse ipv4 address");
			}

			for i := 0; i < 4; i++ {
				x,_ := strconv.Atoi(match[4-i])
				addr = (addr << 8) | uint32(x)
			}

			x,_ := strconv.Atoi(match[5])
			mask = (1 << x) - 1
			addr = (addr & mask) | (^mask)

			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, addr)
			brd = fmt.Sprintf("%d.%d.%d.%d", bs[0], bs[1], bs[2], bs[3])
			break
		}
	}

	return brd
}

func discoverServer() {
	host,port,err := net.SplitHostPort(cluster_listen)
	if err != nil {
		log.Fatal(err)
	}
	brd := calcBroadcast(host) + ":" + port
	host =  ":" + port

	// we need to listen to all addresses to also received broadcast messages
	log.Printf("starting discovery server on %s", host)
	pc,err := net.ListenPacket("udp4", host)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	log.Printf("broadcasting discovery message to %s", brd)
	addr,err := net.ResolveUDPAddr("udp4", brd)
	if err != nil {
		log.Fatal(err)
	}

	msg := writeDiscoverMessage()
	_,err = pc.WriteTo(msg, addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buf := make([]byte, 1024)
		n,addr,err := pc.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		host = addr.String()
		log.Printf("received discovery message from %s", host)
		go readDiscoverMessage(addr.String(), buf[:n])
	}
}

func discoverStart() {
	time.Sleep(time.Second)
	discoverServer()
}
