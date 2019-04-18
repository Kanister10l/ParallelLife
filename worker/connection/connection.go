package connection

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

//ConnectToServer connect to server and handle communication
func ConnectToServer() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/register"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Error dialing server:", err)
		return
	}
	defer c.Close()
	defer log.Println("Closing")

	done := make(chan struct{})

	forwardBoard := make(chan Board)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error trying to read from server:", err)
				return
			}

			log.Println("New task from server")

			board := LoadBoard(string(message))

			forwardBoard <- board
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("Ready"))
	if err != nil {
		log.Println("Error writing message to server:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case t := <-forwardBoard:
			toSend := t.PrepareRetString()

			err := c.WriteMessage(websocket.TextMessage, []byte(toSend))
			if err != nil {
				log.Println("Error writing message to server:", err)
				return
			}

		case <-interrupt:
			defer os.Exit(0)
			log.Println("Closing connection to server")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error writing close message to server:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
