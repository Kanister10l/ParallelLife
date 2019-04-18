package main

import (
	"github.com/kanister10l/ParallelLife/worker/connection"
	"log"
)

func main() {
	log.Println("Worker started")
	connection.ConnectToServer()
}
