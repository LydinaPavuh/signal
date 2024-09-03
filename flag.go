package signal

import (
	"context"
	"github.com/google/uuid"
	"sync"
)

type Flag struct {
	publisher
	raised bool
	mu     sync.Mutex
}

func NewFlag() *Flag {
	return &Flag{publisher: publisher{make(map[uuid.UUID]*Waiter)}}
}

func (fl *Flag) Subscribe() *Waiter {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	waiter := &Waiter{ch: make(chan struct{}, 1)}
	cancel := fl.subscribe(waiter)

	if fl.raised {
		waiter.forceSend()
	}

	waiter.cancel = func() {
		fl.mu.Lock()
		defer fl.mu.Unlock()
		cancel()
	}

	return waiter
}

func (fl *Flag) Raise(ctx context.Context) error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if fl.raised {
		return nil
	}

	fl.raised = true
	return fl.send(ctx, true)
}

func (fl *Flag) Reset() {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fl.raised = false
	fl.purge()
}

func (fl *Flag) IsRaised() bool {
	return fl.raised
}
