package connection

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Board struct {
	Data     []byte
	nextStep []byte
	X        int
	Y        int
	Ret1     int
	Ret2     int
	Mutex    sync.Mutex
}

func LoadBoard(data string) Board {
	split := strings.Split(data, "|")
	if len(split) < 5 {
		log.Println("PANIC ERROR!!! Data reached point of no return. Missing board data. Exit!")
		os.Exit(127)
	}

	x, _ := strconv.Atoi(split[0])
	y, _ := strconv.Atoi(split[1])
	ret1, _ := strconv.Atoi(split[2])
	ret2, _ := strconv.Atoi(split[3])
	boardData := []byte(split[4])

	board := Board{
		X:     x,
		Y:     y,
		Ret1:  ret1,
		Ret2:  ret2,
		Data:  boardData,
		Mutex: sync.Mutex{},
	}

	return board
}

func (b *Board) PrepareRetString() string {

	toSend := fmt.Sprintf("%d|%d|", b.Ret1, b.Ret2)

	for i := 1; i < b.Y-1; i++ {
		for j := 1; j < b.X-1; j++ {
			toSend += string(b.Data[i*b.X+j])
		}
	}

	return toSend
}

func (b *Board) calculateNextBoard() {
	nextStep := make([]byte, b.X*b.Y)
	b.nextStep = nextStep
	end := make(chan bool, b.Y-2)

	for i := 1; i < b.Y-1; i++ {
		go func(endChannel chan bool, b *Board, i int) {
			for j := 1; j < b.X-1; j++ {
				if b.Data[i*b.X+j] == 0 {
					ln := 0
					if b.Data[(i+1)*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[(i+1)*b.X+j] == 1 {
						ln++
					}
					if b.Data[(i+1)*b.X+j-1] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j-1] == 1 {
						ln++
					}
					if b.Data[i*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[i*b.X+j-1] == 1 {
						ln++
					}

					if ln == 3 {
						b.Mutex.Lock()
						b.nextStep[i*b.X+j] = 1
						b.Mutex.Unlock()
					}
				} else {
					ln := 0
					if b.Data[(i+1)*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[(i+1)*b.X+j] == 1 {
						ln++
					}
					if b.Data[(i+1)*b.X+j-1] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j] == 1 {
						ln++
					}
					if b.Data[(i-1)*b.X+j-1] == 1 {
						ln++
					}
					if b.Data[i*b.X+j+1] == 1 {
						ln++
					}
					if b.Data[i*b.X+j-1] == 1 {
						ln++
					}

					if ln == 2 || ln == 3 {
						b.Mutex.Lock()
						b.nextStep[i*b.X+j] = 1
						b.Mutex.Unlock()
					}
				}
			}
			endChannel <- true
		}(end, b, i)
	}

	for i := 1; i < b.Y-1; i++ {
		<- end
	}

	b.Data = b.nextStep
}
