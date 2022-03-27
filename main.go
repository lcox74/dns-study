package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lcox74/dns-study/src/models"
)

func main() {
	fmt.Println("Hello, World!")

	pc, err := net.ListenPacket("udp", ":53")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 2048)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	go doMessage(buf)

	raddr, err := net.ResolveUDPAddr("udp", "1.1.1.1:53")
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}
	defer conn.Close()

	buffer := make([]byte, 2048)

	conn.Write(buf)
	conn.SetReadDeadline(time.Now().Add(1000))

	nlen, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return
	}

	pc.WriteTo(buffer[:nlen], addr)
}

func doMessage(buf []byte) {
	res, _ := models.MarshalDNS(buf)
	resStr, err := json.Marshal(res)

	if err != nil {
		println(err)
	}
	println(string(resStr))
}
