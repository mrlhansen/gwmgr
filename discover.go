package main

import (
	"log"
	"fmt"
	"net"
	"regexp"
	"strings"
	"strconv"
	"encoding/json"
	"encoding/binary"
)

type DiscoverMessage struct {
	WebAddress  string `json:"webAddress"`
	ClusterToken string `json:"clusterToken"`
}

func writeDiscoverMessage() []byte {
	message := &DiscoverMessage{
		WebAddress: webListen,
		ClusterToken: "hej",
	}

	data, _ := json.Marshal(message)
	return data
}

func readDiscoverMessage(addr string, data []byte) {
	var msg DiscoverMessage
	json.Unmarshal(data, &msg)

	host,port,_ := net.SplitHostPort(msg.WebAddress)
	if host == "" || host == "0.0.0.0" {
		host,_,_ = net.SplitHostPort(addr)
	}

	log.Println(host,port)
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
	host,port,err := net.SplitHostPort(clusterListen)
	if err != nil {
		log.Fatal(err)
	}
	brd := calcBroadcast(host) + ":" + port


	log.Printf("starting discovery server on %s", clusterListen)
	pc,err := net.ListenPacket("udp4", clusterListen)
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
	discoverServer()
}
