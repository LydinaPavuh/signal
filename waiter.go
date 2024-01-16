package signal

import "sync"

type Waiter struct {
	ch     chan struct{}
	cancel func()
	mu     sync.Mutex
}

func (w *Waiter) Wait() <-chan struct{} {
	return w.ch
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

func (w *Waiter) send() {
	w.ch <- struct{}{}
}
