package raft

import (
	"net"
	"log"
	"fmt"
	"encoding/gob"
	"time"
)

type DiscoveryMessage struct {
	Version int
	Host string
	Port string
}

type Discovery struct {
	Services map[string]DiscoveryMessage
}

func NewMessage(host string, port string) *DiscoveryMessage {
	return &DiscoveryMessage{
		Version:1,
		Host:host,
		Port:port,
	}
}

func NewServiceDiscovery() *Discovery {
	return &Discovery{
		Services: make(map[string]DiscoveryMessage),
	}
}

func (d *Discovery) handleMessage(msg *DiscoveryMessage) {
	hashKey := fmt.Sprintf("%s:%s", msg.Host, msg.Port)
	if _, ok := d.Services[hashKey]; !ok {
		// the service is not yet registered
		d.Services[hashKey] = *msg
		log.Printf("Registered service: %s\r\n", hashKey)
	}
}

func (d *Discovery) ping(host string, port string) {
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.1:11337")

	if err != nil {
		log.Fatalln("Cannot resolve UDP address", err)
	}

	c, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatalln("Cannot connect to service discovery", err)
	}

	msg := NewMessage(host, port)

	enc := gob.NewEncoder(c)

	enc.Encode(msg)
}

func (d *Discovery) startMulticast() {
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.1:11337")

	if err != nil {
		log.Fatalln("Cannot resolve UDP address 224.0.0.1", err)
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)

	if err != nil {
		log.Fatalf("Cannot start service discovery: %v\r\n", err)
	}

	l.SetReadBuffer(8192)

	go func(c *net.UDPConn) {

		log.Println("Started service discovery", c.LocalAddr())
		for {
			dec := gob.NewDecoder(c)
			msg := &DiscoveryMessage{}
			dec.Decode(msg)
			d.handleMessage(msg)
		}

	}(l)
}

func (d *Discovery) StartDiscovery(host string, port string) {

	d.startMulticast()

	go func() {
		for {
			// send ping every 10sec
			d.ping(host, port)
			time.Sleep(10 * time.Second)
		}
	} ()

}