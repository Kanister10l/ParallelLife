package spinner

import (
	"fmt"
	"time"
)

//Spinner base type for implementing spinners
type Spinner struct {
	message  string
	sequence Sequence
	interval int
	start    chan bool
	stop     chan bool
	dispose  chan bool
	iterator int
}

/*
Init mandatory spinner initialisation function
Usage:
msg -> Message displayed after spinner
iv -> Interval for updaing spinner in miliseconds
seq -> Sequence of characters for the spinner. Obtained as a constant value or user defined.
*/
func (s *Spinner) Init(msg string, iv int, seq Sequence) {
	s.message = msg
	s.sequence = seq
	s.interval = iv
	s.start = make(chan bool)
	s.stop = make(chan bool)
	s.dispose = make(chan bool)
	s.iterator = 0
	go s.run()
}

//Start spinner
func (s *Spinner) Start() {
	go func() {
		s.start <- true
	}()
}

//StartAndWait spinner and wait until it is started (Possibly blocking operation)
func (s *Spinner) StartAndWait() {
	s.start <- true
}

//Pause spinner
func (s *Spinner) Pause() {
	go func() {
		s.stop <- true
	}()
}

//PauseAndWait spinner and wait until it is paused (Possibly blocking operation)
func (s *Spinner) PauseAndWait() {
	s.stop <- true
}

//Stop spinner
func (s *Spinner) Stop() {
	go func() {
		s.dispose <- true
	}()
}

//StopAndWait spinner and wait until it is stopped (Possibly blocking operation)
func (s *Spinner) StopAndWait() {
	s.dispose <- true
}

func (s *Spinner) run() {
	l := 0
	for _, v := range s.sequence {
		if len(v) > l {
			l = len(v)
		}
	}

	select {
	case <-s.start:
		for {
			select {
			case <-s.stop:
			stop:
				for {
					select {
					case <-s.start:
						break stop
					case <-s.dispose:
						fmt.Printf("\n")
						return
					}
				}
			case <-s.dispose:
				fmt.Printf("\n")
				return
			default:
				s.print(l)
				time.Sleep(time.Duration(s.interval) * time.Millisecond)
			}
		}
	case <-s.dispose:
		fmt.Printf("\n")
		return
	}
}

func (s *Spinner) print(maxLength int) {
	toWrite := s.sequence[s.iterator]
	if len(s.sequence[s.iterator]) < maxLength {
		for i := 0; i < maxLength-len(s.sequence[s.iterator]); i++ {
			toWrite += " "
		}
	}

	fmt.Printf("\r%s %s", toWrite, s.message)

	s.iterator++
	if s.iterator == len(s.sequence) {
		s.iterator = 0
	}
}
