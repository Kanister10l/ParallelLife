package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kanister10l/ParallelLife/worker/connection"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func(it chan os.Signal) {
		<-it
		os.Exit(0)
	}(interrupt)

	ip := flag.String("ip", "127.0.0.1", "IP address of the server")
	port := flag.String("port", "8080", "Port on which server is listening")

	flag.Parse()

	log.Println("Worker started")
	for {
		connection.ConnectToServer(*ip, *port)
		time.Sleep(500 * time.Millisecond)
	}
}
