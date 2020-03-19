package server

import (
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
)

type Message struct {
	Body   []byte
}

// Callback is a function for message handling
type Callback func(net.Conn, Message)

type Server struct {
	host string
	port uint64

	addr net.Addr
}

// NewServer creates a new server
func NewServer(host string, port uint64) *Server {
	return &Server{host: host, port: port}
}

// Listen listens for incoming connections
func (s *Server) Listen(c Callback) error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%v", s.host, s.port))
	if err != nil {
		return errors.Wrap(err, "error listening")
	}

	defer listen.Close()

	s.addr = listen.Addr()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go func() {
			err := s.handle(conn, c)
			if err != nil {
				conn.Close()
				log.Print(err)
			}
		}()
	}
}

// Addr returns the endpoint network address.
func (s Server) Addr() net.Addr {
	return s.addr
}

// handle handles incoming requests
func (s *Server) handle(conn net.Conn, c Callback) error {
	for {
		buf, err := Read(conn)
		if err != nil {
			return errors.Wrap(err, "error reading from conn")
		}

		if buf == nil {
			continue
		}

		msg := &Message{Body: buf}

		go c(conn, *msg)
	}
}

// Read reads a message from the connection
func Read(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 1024) // warning, arbitrary buf len

	_, err := conn.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "error reading packet")
	}

	return buf, nil
}

// SendMessage sends an encoded message
func (*Server) SendMessage(conn net.Conn, msg Message) error {
	// defer conn.Close()

	_, err := conn.Write(msg.Body)
	if err != nil {
		return err
	}

	return nil
}
