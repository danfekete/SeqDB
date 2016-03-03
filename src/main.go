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
	"flag"
)

/*
	Setup command line flags
 */

var (
	host = flag.String("h", "localhost:3333", "Which IP and port should the server listen on [ipaddr:port]")
	dbFile = flag.String("d", "seq.db", "The path to the SeqDB database file")
)

const (
	SEQDB_VERSION = "0.1.0"
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
	flag.Parse()

	signalCh := make(chan os.Signal)

	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// TODO: make the database file variable based on nodename
	db, err := bolt.Open(*dbFile, 0600, nil)

	if err != nil {
		log.Fatalf("Cannot connect to database: %v\r\n", err)
	}
	log.Printf("Database file %s opened\r\n", *dbFile)
	seqdb.SetDB(db)

	laddr, err := net.ResolveTCPAddr(CONN_TYPE, *host)

	l, err := net.ListenTCP(CONN_TYPE, laddr)
	log.Printf("Listening on %v\r\n", laddr)

	service := seqdb.NewService()
	go service.Serve(l)

	<-signalCh

	db.Close()
	service.Stop()

	log.Println("Terminating SeqDB")
}
