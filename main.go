package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{net.IPv4(0, 0, 0, 0), 9999, ""})
	if err != nil {
		log.Fatalf(":Cannot listen on UDP 9999: %v", err)
	}

	log.Printf("Listening...")
	buf := make([]byte, 10240)
	for {
		len, from, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP recv error: %v", err)
			continue
		}
		fmt.Printf("RX %s: %q\n", from.String(), buf[:len])
		n, err := udpConn.WriteToUDP([]byte("OK\n"), from)
		switch {
		case err != nil:
			log.Printf("UDP send error: %v", err)
		case n != 3:
			log.Printf("UDP sent %d instead of 3", n)
		}
	}
}
