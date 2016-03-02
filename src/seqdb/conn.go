package seqdb

import (
	"fmt"
	"net"
	"strings"
	"github.com/boltdb/bolt"
	"strconv"
	"log"
)

var (
	Db *bolt.DB
	Lock *BucketLocks = NewBucketLock()
)


func parse(s []string) (int,string) {
	switch s[0] {
	case "SET":
		p := SeqPointer{s[1], s[2], Lock}
		p.Lock()
		value, _ := strconv.ParseUint(s[3], 10, 64)
		p.Set(value)
		p.Unlock()
		return 1,"OK\n"

	case "GET":
		p := SeqPointer{s[1], s[2], Lock}
		p.Lock()
		v := p.Get()
		p.Unlock()
		return 1,fmt.Sprintf("%d\n", v)

	case "INC":
		p := SeqPointer{s[1], s[2], Lock}
		p.Lock()
		v := p.Inc()
		p.Unlock()
		return 1,fmt.Sprintf("%d\n", v)

	case "QUIT":
		return -1,"OK\n"
	}

	return 0,"INVALID COMMAND\n"
}

func SetDB(db *bolt.DB) {
	Db = db
}

func Handle(c net.Conn) error {
	log.Printf("Connection established %s -> %s", c.RemoteAddr(), c.LocalAddr())
	defer c.Close()

	for {

		buf := make([]byte, 65536)
		n, err := c.Read(buf)

		if err != nil { return err }

		s := strings.Trim(string(buf[:n]), " \r\n")
		commands := strings.Split(s, " ")

		ret,response := parse(commands)

		c.Write([]byte(response))

		if ret == -1 {break}
	}

	log.Println("Connection closed", c.RemoteAddr())
	return nil
}
