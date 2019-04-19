package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//ConnectWorker handler for connecting and dispatching job for workers
func ConnectWorker(manager *Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
		if err != nil {
			log.Println("Error upgrading request to websocket from", r.RemoteAddr, "\nError:", err)
			return
		}
		defer conn.Close()
		defer log.Println("Conn closed")

		inlet := make(chan string)
		outlet := make(chan string)
		close := make(chan bool)

		manager.NewWorkerChannel <- Worker{
			InChannel:  inlet,
			OutChannel: outlet,
			Close:      close,
		}

		mt, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from worker:", err)
			return
		}
		log.Println("Ready message received")

		for {
			toSend, ok := <-outlet
			if !ok {
				return
			}
			err = conn.WriteMessage(mt, []byte(toSend))
			if err != nil {
				log.Println("Error writing message to worker:", err)
				break
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from worker:", err)
				break
			}
			inlet <- string(message)
		}
	}
}
