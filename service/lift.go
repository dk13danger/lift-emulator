package service

import (
	"fmt"
	"log"
	"time"
)

// Lift states.
const (
	LiftStateStop = iota
	LiftStateMove
)

// Info represents an structure for communicating with external world.
type Info struct {
	State    int
	Floor    int
	IsMoveUp bool
}

// Lift represents an lift model.
type Lift struct {
	doorsDelay time.Duration
	floorDelay time.Duration
	infoCh     chan *Info
	lockCh     chan struct{}
}

// NewLift creates lift model.
func NewLift(count, height, speed uint, delay time.Duration) (*Lift, error) {
	if speed == 0 {
		return nil, fmt.Errorf("speed can't be zero")
	}
	if height == 0 {
		return nil, fmt.Errorf("floor height must be > 0")
	}
	if count < 5 || count > 10 {
		return nil, fmt.Errorf("floor count must be in range: from 5 to 10")
	}
	if speed > height {
		return nil, fmt.Errorf("Speed can't be more than height, It's super lift")
	}

	return &Lift{
		doorsDelay: delay,
		floorDelay: time.Second * (time.Duration(height / speed)),
		infoCh:     make(chan *Info),
		lockCh:     make(chan struct{}),
	}, nil
}

// Run runs lift. Get outside channels for communicating.
func (l *Lift) Run() (<-chan *Info, chan<- struct{}) {
	return l.infoCh, l.lockCh
}

// Move moves lift from one floor to other.
func (l *Lift) Move(from, to int) {
	if from == to {
		l.OpenCloseDoors()
		return
	}

	isMoveUp := from < to
	if isMoveUp {
		for i := from; i < to; i++ {
			time.Sleep(l.floorDelay)
			log.Printf("Lift move the floor: #%d\n", i)
			l.sendInfo(LiftStateMove, i, isMoveUp)
		}
	} else {
		for i := from; i > to; i-- {
			time.Sleep(l.floorDelay)
			log.Printf("Lift move the floor: #%d\n", i)
			l.sendInfo(LiftStateMove, i, isMoveUp)
		}
	}

	l.sendInfo(LiftStateStop, to, isMoveUp)
}

// OpenCloseDoors opens and closes lift doors by delay time.
func (l *Lift) OpenCloseDoors() {
	log.Printf("Lift open doors\n")
	time.Sleep(l.doorsDelay)
	log.Printf("Lift close doors\n")
}

// sendInfo sends info by channel for dispatcher goroutine.
func (l *Lift) sendInfo(state, floor int, isMoveUp bool) {
	l.infoCh <- &Info{
		State:    state,
		Floor:    floor,
		IsMoveUp: isMoveUp,
	}

	// wait until information will be processed by dispatcher service
	<-l.lockCh
}
