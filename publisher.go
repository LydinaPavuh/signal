package signal

import (
	"context"
	"github.com/google/uuid"
)

type publisher struct {
	subscribers map[uuid.UUID]*Waiter
}

func (sig *publisher) send(ctx context.Context, nowait bool) error {
	for _, w := range sig.subscribers {
		if err := w.send(ctx, nowait); err != nil {
			return err
		}
	}
	return nil
}

func (sig *publisher) purge() {
	for _, w := range sig.subscribers {
		w.Purge()
	}
}

func (sig *publisher) subscribe(w *Waiter) func() {
	key := uuid.Must(uuid.NewRandom())
	sig.subscribers[key] = w
	return func() { delete(sig.subscribers, key) }

}
