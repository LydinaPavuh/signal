package signal

import (
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

func (sig *Flag) Subscribe() *Waiter {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	waiter := &Waiter{ch: make(chan struct{}, 1)}
	cancel := sig.subscribe(waiter)
	waiter.cancel = func() {
		sig.mu.Lock()
		defer sig.mu.Unlock()
		cancel()
	}
	return waiter
}

func (sig *Flag) Raise() {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	if sig.raised {
		return
	}
	sig.raised = true
	sig.send()
}

func (sig *Flag) Reset() {
	sig.mu.Lock()
	defer sig.mu.Unlock()
	sig.raised = false
	sig.purge()
}

func (sig *Flag) IsRaised() bool {
	return sig.raised
}
