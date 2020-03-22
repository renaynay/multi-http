package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Callback func(conn *websocket.Conn) error

func main() {
	Listen()
}

func Listen() {

	http.Handle("/", RouteWSOrHTTP())

	http.ListenAndServe(":8080", nil)
}

func IsWebsocket(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && strings.ToLower(r.Header.Get("Connection")) == "upgrade"
}

func RouteWSOrHTTP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsWebsocket(r) {
			HandleWebsocket(w, r)
			return
		}

		// call HTTP-specific handler

		w.Write([]byte("HTTP!"))
	})
}


func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	// upgrader is needed to upgrade the HTTP Connection to a websocket Connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//Upgrading HTTP Connection to websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}

	for {
		_, raw, err := wsConn.ReadMessage()
		if err != nil {
			log.Error("error reading message from ws connection: ", err)
			return
		}

		log.Info("message received: ", string(raw))

		err = wsConn.WriteMessage(websocket.TextMessage, ([]byte("message received on websocket connection!")))
		if err != nil {
			log.Error("error writing message to ws connection: ", err)
			return
		}
	}
}
