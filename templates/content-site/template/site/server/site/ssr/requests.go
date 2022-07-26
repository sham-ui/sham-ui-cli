package ssr

import (
	"sync"
)

type requests struct {
	sync.Mutex
	storage map[string]chan<- *nodejsResponse
	ch      chan bool
}

func (r *requests) subscribe(msg *nodejsRequest) <-chan *nodejsResponse {
	ch := make(chan *nodejsResponse, 1)
	r.Lock()
	r.storage[msg.ID] = ch
	r.Unlock()
	r.ch <- true
	return ch
}

func (r *requests) publish(msg *nodejsResponse) {
	r.Lock()
	ch, ok := r.storage[msg.ID]
	r.Unlock()
	if ok {
		r.unsubscribe(msg.ID)
		ch <- msg
	}
}

func (r *requests) unsubscribe(id string) {
	r.Lock()
	delete(r.storage, id)
	r.Unlock()
}

func (r *requests) Next() bool {
	return <-r.ch
}

func newRequests(poolSize int) *requests {
	return &requests{
		Mutex:   sync.Mutex{},
		storage: make(map[string]chan<- *nodejsResponse),
		ch:      make(chan bool, poolSize),
	}
}
