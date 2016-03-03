package raft

import (
	"net"
	"log"
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

func (d *Discovery) handleMessage(msg *DiscoveryMessage) {
	
}

func (d *Discovery) StartMulticast() {
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.1")

	if err != nil {
		log.Fatalln("Cannot resolve UDP address 224.0.0.1")
	}

	l, err := net.ListenMulticastUDP("udp", nil, addr)

	if err != nil {
		log.Fatalf("Cannot start service discovery")
	}

	l.SetReadBuffer(8192)


}