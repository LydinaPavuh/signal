package signal

import (
	"context"
	"github.com/google/uuid"
	"sync"
)

type Signal struct {
	publisher
	mu         sync.Mutex
	bufferSize int
}

func NewSignal(bufferSize int) *Signal {
	return &Signal{
		publisher:  publisher{make(map[uuid.UUID]*Waiter)},
		bufferSize: bufferSize,
	}
}

func (sig *Signal) Send(ctx context.Context, nowait bool) error {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	return sig.send(ctx, nowait)
}

func (sig *Signal) Purge() {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	sig.purge()
}

func (sig *Signal) Subscribe() *Waiter {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	waiter := &Waiter{ch: make(chan struct{}, sig.bufferSize)}
	cancel := sig.subscribe(waiter)
	waiter.cancel = func() {
		sig.mu.Lock()
		defer sig.mu.Unlock()
		cancel()
	}
	return waiter
}
