package service

import (
	"sync"
)

// Service represents an service structure.
type Service struct {
	lift        *Lift
	registry    *Registry
	lock        sync.Mutex
	curLiftInfo *Info
	curEvent    *Event
	eventsCh    chan *Event
	liftInfoCh  <-chan *Info
	liftLockCh  chan<- struct{}
}

// New creates new service.
func New(lift *Lift, registry *Registry) *Service {
	initialInfo := &Info{
		Floor:    1,
		IsMoveUp: true,
		State:    LiftStateStop,
	}
	return &Service{
		lift:        lift,
		registry:    registry,
		lock:        sync.Mutex{},
		curLiftInfo: initialInfo,
		eventsCh:    make(chan *Event, 100),
	}
}

// Run runs service.
func (s *Service) Run() chan<- *Event {
	s.liftInfoCh, s.liftLockCh = s.lift.Run()

	go s.processLift()
	go s.processEvents()

	return s.eventsCh
}

// processLift processes all lift info states.
func (s *Service) processLift() {
	for info := range s.liftInfoCh {
		s.lock.Lock()
		s.curLiftInfo = info
		s.lock.Unlock()

		switch info.State {
		case LiftStateStop:
			s.lift.OpenCloseDoors()
			s.registry.Drop(s.curEvent)

			// find any next "external" event from registry
			if e := s.registry.GetAnyDiffer(ExternalCall, info.Floor); e != nil {
				// run lift immediately for next floor which "waiting" lift
				s.moveLift(e)
				break
			}

			// if has no next external events and lift has stopped then clear all internal calls
			// because this internal calls was useful only for the "moving lift session"
			events := s.registry.GetByType(InternalCall)
			for _, e := range events {
				s.registry.Drop(e)
			}

		case LiftStateMove:
			// if lift move up - process only internal calls
			if info.IsMoveUp {
				if e := s.registry.GetEqual(InternalCall, info.Floor); e != nil {
					s.lift.OpenCloseDoors()
					s.registry.Drop(s.curEvent)
				}
				break
			}

			// if lift move down - process internal and external calls (get events only by floor)
			events := s.registry.GetByFloor(info.Floor)
			if len(events) > 0 {
				s.lift.OpenCloseDoors()
				for _, e := range events {
					s.registry.Drop(e)
				}
			}
		}

		// unlock lift. If the lift has "moving state" then lift continue moving
		s.liftLockCh <- struct{}{}
	}
}

// processEvents process events from web server, parses it and put into internal cache (events registry).
func (s *Service) processEvents() {
	for e := range s.eventsCh {
		s.registry.Add(e)
		if e.Type == ExternalCall {
			s.moveLift(e)
		}
	}
}

// moveLift moves lift async (in separate goroutine). Thread safe.
func (s *Service) moveLift(e *Event) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.curLiftInfo.State != LiftStateStop {
		return
	}

	s.curLiftInfo.State = LiftStateMove
	s.curEvent = e

	go s.lift.Move(s.curLiftInfo.Floor, e.Floor)
}
