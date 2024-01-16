package signal

import (
	"github.com/google/uuid"
)

type publisher struct {
	subscribers map[uuid.UUID]*Waiter
}

func (sig *publisher) send() {
	for _, w := range sig.subscribers {
		w.send()
	}
}

func (sig *publisher) purge() {
	for _, w := range sig.subscribers {
		w.Purge()
	}
}

func (sig *publisher) subscribe(w *Waiter) func() {
	key := uuid.Must(uuid.NewRandom())
	sig.subscribers[key] = w
	return func() {
		delete(sig.subscribers, key)
	}
}
