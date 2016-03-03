package seqdb

import (
	"sync"
	"log"
	"net"
	"time"
	"strings"
)

type Service struct {
	ch chan bool
	waitGroup *sync.WaitGroup
}

// Make a new Service.
func NewService() *Service {
	s := &Service{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	s.waitGroup.Add(1)
	return s
}

// Accept connections and spawn a goroutine to serve each one.  Stop listening
// if anything is received on the service's channel.
func (s *Service) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("stopping listening on", listener.Addr())
			listener.Close()
			return
		default:
		}

		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()

		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}

		log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)

		go s.serve(conn)
	}
}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (s *Service) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}

// Serve a connection by reading and writing what was read.  That's right, this
// is an echo service.  Stop reading and writing if anything is received on the
// service's channel but only after writing what was read.
func (s *Service) serve(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	for {

		select {
		case <-s.ch:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}

		conn.SetDeadline(time.Now().Add(10*1e9))

		buf := make([]byte, 4096)

		n, err := conn.Read(buf)

		if err != nil {
			log.Printf("Cannot read from buffer: %v\r\n", err)

		}

		s := strings.Trim(string(buf[:n]), " \r\n")

		// Parse the command
		code, response := ParseCommand(s)

		// QUIT command was issued
		if code == -1 { break }

		// Write response
		_, err = conn.Write([]byte(response))

		if err != nil {
			log.Printf("Error writing buffer: %v\r\n", err)
		}
	}
}