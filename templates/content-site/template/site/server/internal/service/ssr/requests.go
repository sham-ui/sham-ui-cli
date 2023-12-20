package ssr

import (
	"context"
	"site/pkg/tracing"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

type requestHolder struct {
	span       trace.Span
	subscriber chan *nodejsResponse
}

func (rh *requestHolder) start(parentCtx context.Context) {
	const op = scopeName + "requests.processing"
	ps := tracing.SpanFromContext(parentCtx)
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), ps.SpanContext())
	_, rh.span = ps.TracerProvider().Tracer(scopeName).Start(ctx, op)
}

func (rh *requestHolder) end() {
	rh.span.End()
}

func newRequestHolder() *requestHolder {
	return &requestHolder{
		span:       nil,
		subscriber: make(chan *nodejsResponse, 1),
	}
}

type requests struct {
	sync.Mutex
	storage map[string]*requestHolder
	ch      chan bool
}

func (r *requests) subscribe(ctx context.Context, msg *nodejsRequest) <-chan *nodejsResponse {
	holder := newRequestHolder()

	r.Lock()
	r.storage[msg.ID] = holder
	r.Unlock()

	holder.start(ctx)

	r.ch <- true
	return holder.subscriber
}

func (r *requests) publish(msg *nodejsResponse) {
	r.Lock()
	holder, exists := r.storage[msg.ID]
	if exists {
		delete(r.storage, msg.ID)
	}
	r.Unlock()

	if exists {
		holder.end()
		holder.subscriber <- msg
	}
}
func (r *requests) unsubscribe(id string) {
	r.Lock()
	holder, exists := r.storage[id]
	r.Unlock()
	if !exists {
		return
	}
	delete(r.storage, id)
	holder.end()
}

func (r *requests) Next() bool {
	return <-r.ch
}

func newRequests(poolSize int) *requests {
	return &requests{
		Mutex:   sync.Mutex{},
		storage: make(map[string]*requestHolder),
		ch:      make(chan bool, poolSize),
	}
}
