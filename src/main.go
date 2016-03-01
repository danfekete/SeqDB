package main

import (
	"fmt"
	"net"
	"os"
	"github.com/boltdb/bolt"
	//"encoding/json"
	"seqdb"
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
		fmt.Println("Cannot connect to database!")
		os.Exit(1)
	}

	seqdb.SetDB(db)

	l, err := net.Listen(CONN_TYPE, "localhost:3333")
	if err != nil {
		fmt.Println("Cannot listen on port")
		os.Exit(1)
	}

	defer l.Close()
	defer db.Close()

	for {
		// Yummy infinte loop

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go seqdb.Handle(conn)
	}
}
