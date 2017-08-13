package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	udpConn, err := net.ListenUDP("udp4", &net.UDPAddr{net.IPv4(0, 0, 0, 0), 9999, ""})
	if err != nil {
		log.Fatalf(":Cannot listen on UDP 9999: %v", err)
	}

	log.Printf("Listening to UDP 9999...")
	for {
		buf := make([]byte, 4096)
		len, from, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP recv error: %v", err)
			continue
		}
		if len > 2 && buf[0] == '{' {
			go func() {
				if err := handleJSON(udpConn, from, buf[:len]); err != nil {
					log.Println(err)
				}
			}()
		} else {
			fmt.Printf("RX %s: %q\n", from.String(), buf[:len])
			n, err := udpConn.WriteToUDP([]byte("??\n"), from)
			switch {
			case err != nil:
				log.Printf("UDP send error: %v", err)
			case n != 3:
				log.Printf("UDP sent %d instead of 3", n)
			}
		}
	}
}

type message struct {
	Ack  int    `json:"a"`
	Key  string `json:"k"`
	Data string `json:"d"`
	Tags string `json:"t"`
}

func handleJSON(conn *net.UDPConn, from *net.UDPAddr, buf []byte) error {
	var msg message
	if err := json.Unmarshal(buf, &msg); err != nil {
		return fmt.Errorf("cannot parse json message: %v", err)
	}
	fmt.Printf("RX %s: JSON %+v\n", from.String(), msg)

	tcpConn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{net.IPv4(23, 253, 146, 203), 9999, ""})
	if err != nil {
		conn.WriteToUDP([]byte("[7,0]"), from)
		return fmt.Errorf("cannot connect to Hologram: %v", err)
	}
	tcpConn.SetDeadline(time.Now().Add(25 * time.Second))
	tcpConn.Write(buf)
	var resp [128]byte
	n, err := tcpConn.Read(resp[:])
	if err != nil {
		conn.WriteToUDP([]byte("[7,0]"), from)
		return fmt.Errorf("cannot get response from Hologram: %v", err)
	}
	fmt.Printf("HO: %q\n", resp[:n])

	if msg.Ack > 0 && msg.Ack < 3610 {
		fmt.Printf("Delay %d seconds\n", msg.Ack)
		time.Sleep(time.Duration(msg.Ack) * time.Second)
	}
	n, err = conn.WriteToUDP(resp[:n], from)
	switch {
	case err != nil:
		return fmt.Errorf("UDP send error: %v", err)
	case n != 3:
		return fmt.Errorf("UDP sent %d instead of 3", n)
	}
	return nil
}
