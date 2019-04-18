package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/kanister10l/ParallelLife/server"
)

func main() {
	logLevel := flag.String("ll", "clean", "Log Level. Available levels clean|debug")

	flag.Parse()

	if *logLevel == "clean" {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	game := server.NewGame(100, 100, 0.1)
	server.NewManager(game)

	mux := bone.New()

	mux.Handle("/register", server.ConnectWorker())

	log.Println("Listening on 0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", mux)
}
