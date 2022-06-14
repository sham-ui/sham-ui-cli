package ssr

import (
	"sync"
)

type requests struct {
	sync.Mutex
	storage map[string]chan<- *nodejsResponse
}

func (r *requests) subscribe(msg *nodejsRequest, ch chan<- *nodejsResponse) {
	r.Lock()
	r.storage[msg.ID] = ch
	r.Unlock()
}

func (r *requests) publish(msg *nodejsResponse) {
	r.Lock()
	ch, ok := r.storage[msg.ID]
	if ok {
		delete(r.storage, msg.ID)
		r.Unlock()
		ch <- msg
	} else {
		r.Unlock()
	}
}

func newRequests() *requests {
	return &requests{
		Mutex:   sync.Mutex{},
		storage: make(map[string]chan<- *nodejsResponse),
	}
}
