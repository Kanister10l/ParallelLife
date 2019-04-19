package server

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kanister10l/ParallelLife/spinner"
)

//Manager base type required for game manager to work
type Manager struct {
	Game             *Game
	Workers          []Worker
	NewWorkerChannel chan Worker
	Mutex            sync.Mutex
	GlueChannel      chan string
	Next             chan bool
	Gif              chan gifBoard
	Generations      int
}

//Worker contains channel for communication with connected worker
type Worker struct {
	InChannel  chan string
	OutChannel chan string
	Close      chan bool
}

//NewManager creates new game manager
func NewManager(game *Game, gens, gifScale, gifDelay int, gifFile string) *Manager {
	manager := &Manager{}
	manager.Game = game
	manager.Workers = []Worker{}
	manager.NewWorkerChannel = make(chan Worker, 100)
	manager.GlueChannel = make(chan string, 100)
	manager.Next = make(chan bool)
	manager.Gif = make(chan gifBoard, gens)
	manager.Generations = gens
	go manager.listenForNewWorker()
	go manager.dispatchJobs()
	go manager.glueBoard()
	go CreateGif(gifFile, game.X, game.Y, gifScale, gifDelay, manager.Gif)

	return manager
}

func (m *Manager) listenForNewWorker() {
	for w := range m.NewWorkerChannel {
		m.Mutex.Lock()
		m.Workers = append(m.Workers, w)
		m.Mutex.Unlock()
	}
}

func (m *Manager) dispatchJobs() {
	spin := spinner.Spinner{}
	spin.Init("Worker warmup", 70, spinner.Circle1())
	spin.StartAndWait()
	time.Sleep(2 * time.Second)
	spin.StopAndWait()

	spin = spinner.Spinner{}
	spin.Init("Simulating", 70, spinner.Circle1())
	spin.StartAndWait()

	for {
		m.Mutex.Lock()
		workerData := m.Game.PrepareStrings(len(m.Workers))
		for k := range m.Workers {
			m.Workers[k].OutChannel <- workerData[k]
			go func(glue chan string, inlet chan string) {
				select {
				case data := <-inlet:
					glue <- data
				}
			}(m.GlueChannel, m.Workers[k].InChannel)
		}
		m.Mutex.Unlock()

		_, ok := <-m.Next
		if !ok {
			for k := range m.Workers {
				close(m.Workers[k].OutChannel)
			}
			time.Sleep(400 * time.Millisecond)
			spin.StopAndWait()
			os.Exit(0)
			return
		}
	}
}

func (m *Manager) glueBoard() {
	ready := 0
	toFill := m.Game.X * m.Game.Y
	filled := 0
	newBoard := make([]byte, (m.Game.X+2)*(m.Game.Y+2))

	for data := range m.GlueChannel {
		split := strings.Split(data, "|")
		if len(split) < 3 {
			log.Println("PANIC ERROR!!! Data reached point of no return. Missing board data. Exit!")
			os.Exit(127)
		}

		partY1, _ := strconv.Atoi(split[0])
		partY2, _ := strconv.Atoi(split[1])
		filled += (partY2 - partY1) * m.Game.X
		byteData := []byte(split[2])
		iterator := 0

		for i := partY1; i < partY2; i++ {
			for j := 1; j < m.Game.X+1; j++ {
				newBoard[(i+1)*(m.Game.X+2)+j] = byteData[iterator]
				iterator++
			}
		}

		if toFill == filled {
			m.Mutex.Lock()
			m.Gif <- gifBoard{Data: m.Game.Board, X: m.Game.X, Y: m.Game.Y}
			m.Game.Board = newBoard
			m.Mutex.Unlock()
			ready++
			if ready == m.Generations {
				m.Gif <- gifBoard{Data: m.Game.Board, X: m.Game.X, Y: m.Game.Y}
				close(m.Gif)
				close(m.Next)
				return
			}
			newBoard = make([]byte, (m.Game.X+2)*(m.Game.Y+2))
			filled = 0
			m.Next <- true
		}
	}
}
