package main

import (
	"fmt"
	"net"
	"github.com/boltdb/bolt"
	//"encoding/json"
	"seqdb"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	SEQDB_VERSION = "0.0.2"
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)


type Message struct {
	BucketName string
	SequenceName string
	Value uint64
}

func main() {
	fmt.Println("SeqDB v.", SEQDB_VERSION)
	fmt.Println("Written by Daniel Fekete <daniel.fekete@voov.hu>")

	signalCh := make(chan os.Signal)

	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// TODO: make the database file variable based on nodename
	db, err := bolt.Open("seq.db", 0600, nil)

	if err != nil {
		log.Fatalf("Cannot connect to database: %v\r\n", err)
	}

	seqdb.SetDB(db)

	laddr, err := net.ResolveTCPAddr(CONN_TYPE, fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT))

	l, err := net.ListenTCP(CONN_TYPE, laddr)

	service := seqdb.NewService()
	go service.Serve(l)

	<-signalCh

	db.Close()
	service.Stop()

	log.Println("Terminating SeqDB")
}
