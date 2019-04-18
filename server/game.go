package server

import (
	"fmt"
	"math/rand"
	"time"
)

//Game base struct containing data about current state of the Game of Life
type Game struct {
	Board []byte
	Dimensions
}

//Dimensions structure describing dimensions of the game board
type Dimensions struct {
	X int
	Y int
}

//NewGame prepares and returns game struct
func NewGame(x, y int, chance float64) *Game {
	rand.Seed(time.Now().Unix())

	game := &Game{}
	game.X = x
	game.Y = y
	game.Board = make([]byte, (x+2)*(y+2))
	for i := 1; i < x+1; i++ {
		for j := 1; j < y+1; j++ {
			game.Board[j*(x+2)+i] = randCell(chance)
		}
	}

	return game
}

func randCell(p float64) byte {
	exp := rand.Float64()
	if exp < p {
		return 1
	}

	return 0
}

//PrepareStrings divides game board into worker readable data
func (g *Game) PrepareStrings(n int) []string {
	ret := []string{}

	div := g.Y / n
	for i := 0; i < n; i++ {
		if i < n-1 {
			toSend := fmt.Sprintf("%d|%d|%d|%d|", g.X+2, div+2, i*div, (i+1)*div)
			toSend += string(g.Board[i*div*(g.X+2) : (i+1)*(div+2)*(g.X+2)+g.X+2])
			ret = append(ret, toSend)
		} else {
			toSend := fmt.Sprintf("%d|%d|%d|%d|", g.X+2, g.Y-(i*div)+2, i*div, g.Y)
			toSend += string(g.Board[i*div*(g.X+2):])
			ret = append(ret, toSend)
		}
	}

	return ret
}
