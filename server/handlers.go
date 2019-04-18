package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//ConnectWorker handler for connecting and dispatching job for workers
func ConnectWorker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			log.Println("Error upgrading request to websocket from", r.RemoteAddr, "\nError:", err)
			return
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from worker:", err)
				break
			}
			log.Printf("Recived message from worker: %s", message)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("Error writing message to worker:", err)
				break
			}
		}
	}
}
