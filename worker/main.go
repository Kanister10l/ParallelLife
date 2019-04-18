package main

import (
	"log"
	"time"

	"github.com/kanister10l/ParallelLife/worker/connection"
)

func main() {
	log.Println("Worker started")
	for {
		connection.ConnectToServer()
		log.Println("Connection Failed, Retrying after 500ms")
		time.Sleep(500 * time.Millisecond)
	}
}
