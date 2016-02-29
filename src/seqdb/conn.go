package seqdb

import (
	"fmt"
	"net"
	"strings"
	"github.com/boltdb/bolt"
	"strconv"
)

var (
	Db *bolt.DB
	BucketLocks map[string]bool
)


func parse(s []string) (int,string) {
	switch s[0] {
	case "SET":
		p := SeqPointer{s[1], s[2]}
		value, _ := strconv.ParseUint(s[3], 10, 64)
		p.Set(value)
		return 1,"OK\n"
	case "GET":
		p := SeqPointer{s[1], s[2]}
		v := p.Get()
		fmt.Println("Getting", v)
		return 1,fmt.Sprintf("%d\n", v)
	case "INC":
		p := SeqPointer{s[1], s[2]}
		v := p.Inc()
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
	fmt.Println("Connection established", c.RemoteAddr())
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

	fmt.Println("Connection closed", c.RemoteAddr())
	return nil
}
