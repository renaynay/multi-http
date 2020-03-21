package multi

import (
	"net"
	"net/http"

	"github.com/atoulme/aproxi/message"
	log "github.com/sirupsen/logrus"
	"github.com/renaynay/multi-http/server"
	ws "github.com/gorilla/websocket"
)

// 1. listen on endpoint (net.Listen)
// 2. create http.Server and ws.Server ( this is where you'd use upgrader )
	// func (s *Server) WebsocketHandler(allowedOrigins []string) http.Handler {
		//	var upgrader = websocket.Upgrader{
		//		ReadBufferSize:  wsReadBuffer,
		//		WriteBufferSize: wsWriteBuffer,
		//		WriteBufferPool: wsBufferPool,
		//		CheckOrigin:     wsHandshakeValidator(allowedOrigins),
		//	}
		//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		conn, err := upgrader.Upgrade(w, r, nil)
		//		if err != nil {
		//			log.Debug("WebSocket upgrade failed", "err", err)
		//			return
		//		}
		//		codec := newWebsocketCodec(conn)
		//		s.ServeCodec(codec, 0)
		//	})
		//}
// 3. serve on the http.Server (server.Serve(listener)), this just means accept incoming connections on the Listener, creating a
////  new service goroutine for each
// 4.

type HTTPServer struct {
	Server *server.Server
	HTTP *http.Server
}

func NewHTTPServer(srv *server.Server) *HTTPServer {
	return &HTTPServer{
		Server: srv,
		HTTP:   &http.Server{
			Handler: NewWSHandler(),
		},
	}
}

func NewWSHandler() http.Handler {
	upgrader := ws.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Debug("WebSocket upgrade failed", "err", err)
			return
		}
	})
}

func Start(srv *server.Server) {
	wsSrv := NewHTTPServer(srv)
	wsSrv.Server.Listen(func(_ net.Conn, msg message.Message){

	})
}

//func NewHTTPServer(host string, port uint64) *server.Server {
//	srv := server.NewServer(host, port)
//	httpSrv := http.Server{
//		Addr:              "",
//		Handler:           nil,
//		TLSConfig:         nil,
//		ReadTimeout:       0,
//		ReadHeaderTimeout: 0,
//		WriteTimeout:      0,
//		IdleTimeout:       0,
//		MaxHeaderBytes:    0,
//		TLSNextProto:      nil,
//		ConnState:         nil,
//		ErrorLog:          nil,
//		BaseContext:       nil,
//		ConnContext:       nil,
//	}
//	httpSrv.ListenAndServe()
//
//	upgrader := ws.Upgrader{}
//}
