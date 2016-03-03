package seqdb

import (
	"fmt"
	"strings"
	"github.com/boltdb/bolt"
	"strconv"
	"log"
	"regexp"
)

var (
	Db *bolt.DB
	Lock *BucketLocks = NewBucketLock()
)


func parse(s []string) (int,string) {
	switch s[0] {
	case "SET":
		if len(s) < 4 {
			return 0,"BAD ARGUMENTS\n"
		}
		p := SeqPointer{s[1], s[2], Lock}
		p.Lock()
		value, _ := strconv.ParseUint(s[3], 10, 64)
		p.Set(value)
		p.Unlock()
		return 1,"OK\n"

	case "GET":
		if len(s) < 3 {
			return 0,"BAD ARGUMENTS\n"
		}
		p := SeqPointer{s[1], s[2], Lock}
		p.Lock()
		v := p.Get()
		p.Unlock()
		return 1,fmt.Sprintf("%d\n", v)

	case "INC":
		if len(s) < 3 {
			return 0,"BAD ARGUMENTS\n"
		}
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

func ParseCommand(cmd string) (int, string) {
	s := strings.Trim(cmd, " \r\n")
	r, err := regexp.Compile("^([A-Z]{3,6})( [-_A-Za-z0-9]+)*$")

	if err != nil {
		log.Printf("Cannot create regexp: %v\r\n", err)
	}

	// check for command validity
	if r.MatchString(s) == false {
		// the command is malformed
		return 0, "INVALID COMMAND"
	}

	commands := strings.Split(s, " ")
	return parse(commands)
}
