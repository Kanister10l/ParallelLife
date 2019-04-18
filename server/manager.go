package server

import (
	"strconv"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kanister10l/ParallelLife/spinner"
)

type Manager struct {
	Game             *Game
	Workers          []Worker
	NewWorkerChannel chan Worker
	Mutex            sync.Mutex
	GlueChannel      chan string
	Next             chan bool
}

type Worker struct {
	InChannel  chan string
	OutChannel chan string
	Close      chan bool
}

func NewManager(game *Game) *Manager {
	manager := &Manager{}
	manager.Workers = []Worker{}
	manager.NewWorkerChannel = make(chan Worker, 100)
	manager.GlueChannel = make(chan string, 100)
	manager.Next = make(chan bool)
	go manager.listenForNewWorker()
	go manager.dispatchJobs()

	return manager
}

func (m *Manager) listenForNewWorker() {
	for w := range m.NewWorkerChannel {
		m.Mutex.Lock()
		m.Workers = append(m.Workers, w)
		//TODO: Remove worker on close message
		m.Mutex.Unlock()
	}
}

func (m *Manager) dispatchJobs() {
	spin := spinner.Spinner{}
	spin.Init("Worker warmup", 70, spinner.Circle1())
	spin.StartAndWait()
	time.Sleep(5 * time.Second)
	spin.StopAndWait()

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

		<-m.Next
	}
}

func (m *Manager) glueBoard() {
	toFill := m.Game.X * m.Game.Y
	filled := 0
	newBoard := make([]byte, (m.Game.X + 2)*(m.Game.Y + 2))

	for data := range m.GlueChannel {
		split := strings.Split(data, "|")
		if len(split) < 3 {
			log.Println("PANIC ERROR!!! Data reached point of no return. Missing board data. Exit!")
			os.Exit(127)
		}

		partY1,_ := strconv.Atoi(split[0])
		partY2,_ := strconv.Atoi(split[1])
		filled += (partY2 - partY1) * m.Game.X
		byteData := []byte(split[2])
		iterator := 0

		for i := partY1; i < partY2; i++ {
			for j := 1; j < m.Game.X + 1; j++ {
				newBoard[(i + 1) * (m.Game.X + 2) + j] = byteData[iterator]
				iterator++
			}
		}

		if toFill == filled {
			m.Mutex.Lock()
			m.Game.Board = newBoard
			m.Mutex.Unlock()
			newBoard = make([]byte, (m.Game.X + 2)*(m.Game.Y + 2))
			filled = 0
			m.Next <- true
		}
	}
}
