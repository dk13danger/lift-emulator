package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

// Call types (external, internal).
const (
	ExternalCall = iota
	InternalCall
)

// Event represents structure for processing events from web.
type Event struct {
	Type  int
	Floor int
}

// GetID gets id for event. Use MD5 hash.
func (e *Event) GetID() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%d-%d", e.Type, e.Floor)))
	return hex.EncodeToString(h.Sum(nil))
}

// Registry represents events registry.
type Registry struct {
	lock   sync.RWMutex
	events map[string]*Event
}

// NewRegistry creates events registry.
func NewRegistry() *Registry {
	return &Registry{
		lock:   sync.RWMutex{},
		events: make(map[string]*Event),
	}
}

// Add adds event to registry.
func (r *Registry) Add(e *Event) {
	r.lock.Lock()
	defer r.lock.Unlock()

	id := e.GetID()
	if _, ok := r.events[id]; !ok {
		r.events[id] = e
	}
}

// GetAnyDiffer gets any different event from registry.
func (r *Registry) GetAnyDiffer(eventType, floor int) *Event {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, v := range r.events {
		if v.Floor != floor && v.Type == eventType {
			return v
		}
	}
	return nil
}

// GetByFloor gets all events by corresponding floor.
func (r *Registry) GetByFloor(floor int) []*Event {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var events []*Event
	for _, v := range r.events {
		if v.Floor == floor {
			events = append(events, v)
		}
	}
	return events
}

// GetByType gets all events by corresponding type.
func (r *Registry) GetByType(eventType int) []*Event {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var events []*Event
	for _, v := range r.events {
		if v.Type == eventType {
			events = append(events, v)
		}
	}
	return events
}

// GetEqual gets equal event from registry.
func (r *Registry) GetEqual(eventType, floor int) *Event {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, v := range r.events {
		if v.Floor == floor && v.Type == eventType {
			return v
		}
	}
	return nil
}

// Drop drops event from registry
func (r *Registry) Drop(e *Event) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	id := e.GetID()
	if _, ok := r.events[id]; ok {
		delete(r.events, id)
		return true
	}
	return false
}
