package server_test

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"github.com/renaynay/multi-http/server"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	m.Run()
}

func TestTCP_Listen(t *testing.T) {
	srv := server.NewServer("127.0.0.1", uint64(0))
	ch := make(chan server.Message)

	go srv.Listen(func(_ net.Conn, msg server.Message) {
		ch <- msg
	})

	for {
		if srv.Addr() != nil {
			break
		}

		time.Sleep(1)
	}

	conn, err := net.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Error("could not connect to TCP srv: ", err)
	}

	buf := make([]byte, 1014) // 10 less than 1024
	body := append([]byte("testing123"), buf...) // prepend "testing123" to buf = 1024 bytes long

	msg := server.Message{Body: body}

	_, err = conn.Write(msg.Body)
	if err != nil {
		t.Error("could not write to the TCP srv: ", err)
	}
	read := <-ch

	expected := server.Message{Body: body}

	assert.Equal(t, expected, read)
}

func TestTCP_SendMessage(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Error("could not listen on host: ", err)
	}
	defer listener.Close()

	msgs := make(chan server.Message)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error("error accepting connection: ", err)
		}

		buf := make([]byte, 1024)

		_, err = io.ReadFull(conn, buf)
		if err != nil {
			t.Error("error reading from connection: ", err)
		}

		msg := &server.Message{Body: buf}
		if err != nil {
			t.Error("error unmarshaling message: ", err)
		}

		msgs <- *msg
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		t.Error("failed to dial host", err)
	}
	defer conn.Close()

	buf := make([]byte, 1014) // 10 less than 1024
	body := append([]byte("testing123"), buf...) // prepend "testing123" to buf = 1024 bytes long

	srv := server.NewServer("127.0.0.1", 3000)
	srv.SendMessage(conn, server.Message{body})

	msg := <-msgs
	expected := server.Message{Body: body}

	assert.Equal(t, expected, msg)
}