package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/kanister10l/ParallelLife/server"
)

func main() {
	logLevel := flag.String("ll", "clean", "Log Level. Available levels clean|debug")
	gifFile := flag.String("o", "a.gif", "FIlen name of the output gif file")
	x := flag.Int("x", 100, "Horizontal size of game board (Integer)")
	y := flag.Int("y", 100, "Vertical size of game board (Integer)")
	scale := flag.Int("s", 1, "Scale for output gif file (Integer)")
	chance := flag.Float64("pr", 0.1, "Probability for cell to spawn (Float64)")
	gens := flag.Int("g", 50, "Number of generations for simulation (Integer)")
	delay := flag.Int("d", 100, "Delay between gif frames (Integer - increments of 10)")
	ip := flag.String("ip", "0.0.0.0", "IP address of the server")
	port := flag.String("port", "8080", "Port on which server is listening")

	flag.Parse()

	if *logLevel == "clean" {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	log.Println(*gifFile)
	game := server.NewGame(*x, *y, *chance)
	manager := server.NewManager(game, *gens, *scale, (*delay)/10, *gifFile)

	mux := bone.New()

	mux.Get("/register", server.ConnectWorker(manager))

	http.ListenAndServe(fmt.Sprintf("%s:%s", *ip, *port), mux)
}
