package main

import (
	"fmt"
	"net"
	"github.com/boltdb/bolt"
	//"encoding/json"
	"seqdb"
	"log"
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

	db, err := bolt.Open("seq.db", 0600, nil)

	if err != nil {
		log.Fatalf("Cannot connect to database: %v\r\n", err)
	}

	seqdb.SetDB(db)

	l, err := net.Listen(CONN_TYPE, "localhost:3333")
	if err != nil {
		log.Fatalf("Cannot listen on %s:%d: %v\r\n", CONN_HOST, CONN_PORT, err)
	}

	defer l.Close()
	defer db.Close()

	for {
		// Yummy infinte loop

		conn, err := l.Accept()

		if err != nil {
			log.Fatalf("Error when accepting connection %v\r\n", err)
		}

		go seqdb.Handle(conn)
	}
}
