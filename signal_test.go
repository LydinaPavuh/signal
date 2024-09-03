package signal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	s := NewSignal(1)
	waiter := s.Subscribe()
	s.Send(context.Background(), true)
	select {
	case <-waiter.Wait():
	default:
		t.Fatalf("Signal not received")
	}

	s.Send(context.Background(), true)
	s.Purge()

	select {
	case <-waiter.Wait():
		t.Fatalf("Unexpected signal received")
	default:
	}

	// Cancel waiter and try to send signal
	waiter.Cancel()
	s.Send(context.Background(), true)
	if !waiter.Empty() {
		t.Fatalf("Unexpected signal received")
	}
}

func TestSignalConcurrentSubscribe(t *testing.T) {
	concurrentWaiters := 10000

	signal := NewSignal(1)
	gate := NewFlag()
	ch := make(chan *Waiter, concurrentWaiters)

	for i := 0; i < concurrentWaiters; i++ {
		gateWaiter := gate.Subscribe()
		go func() {
			assert.NoError(t, gateWaiter.WaitBlocking(context.Background()))
			ch <- signal.Subscribe()
		}()
	}

	assert.NoError(t, gate.Raise(context.Background()))
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, signal.Send(context.Background(), true))

	var waiter *Waiter
	vc := 0
	for {
		vc++
		select {
		case waiter = <-ch:
			select {
			case <-waiter.Wait():
				t.Logf("Waiter %d raised correct", vc)

			default:
				assert.FailNow(t, "Waiter dont receive signal")
			}

		default:
			return
		}
	}
}
