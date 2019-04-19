package connection

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

//ConnectToServer connect to server and handle communication
func ConnectToServer(ip, port string) {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", ip, port), Path: "/register"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer c.Close()
	defer log.Println("Closing")

	done := make(chan struct{})

	forwardBoard := make(chan *Board)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}

			log.Println("New task from server")

			board := LoadBoard(string(message))
			board.calculateNextBoard()

			forwardBoard <- &board
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
		}
	}
}
