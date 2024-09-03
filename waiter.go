package signal

import (
	"context"
	"sync"
)

type Waiter struct {
	ch     chan struct{}
	cancel func()
	mu     sync.Mutex
}

func (w *Waiter) Wait() <-chan struct{} {
	return w.ch
}

func (w *Waiter) WaitBlocking(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.ch:
		return nil
	}
}

func (w *Waiter) Cancel() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cancel()
	close(w.ch)
}

func (w *Waiter) Empty() bool {
	return len(w.ch) == 0
}

func (w *Waiter) Purge() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.ch) > 0 {
	purgeLoop:
		for {
			select {
			case <-w.ch:
			default:
				break purgeLoop
			}
		}
	}
}

func (w *Waiter) send(ctx context.Context, nowait bool) error {
	if nowait {
		return w.sendNowait(ctx)
	} else {
		return w.sendWait(ctx)
	}
}

func (w *Waiter) sendWait(ctx context.Context) error {
	select {
	case w.ch <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (w *Waiter) sendNowait(ctx context.Context) error {
	select {
	case w.ch <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	default: // do nothing
	}
	return nil
}

func (w *Waiter) forceSend() {
	select {
	case w.ch <- struct{}{}:
	default: // do nothing
	}
}
