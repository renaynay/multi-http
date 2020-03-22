package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("websocket!"))
	}).Headers("Connection", "upgrade", "Upgrade", "websocket", "Sec-Websocket-Version", "13")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("just plain http!"))
	})

	http.ListenAndServe(":8080", router)
}
