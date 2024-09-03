package signal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFlag(t *testing.T) {
	f := NewFlag()
	waiter := f.Subscribe()

	t.Run("Raise", func(t *testing.T) {

		assert.NoError(t, f.Raise(context.Background()))
		assert.True(t, f.IsRaised())

		select {
		case <-waiter.Wait():
		default:
			t.Fatalf("Flag not raised")
		}
	})

	t.Run("Post rise subscribe", func(t *testing.T) {
		waiter := f.Subscribe()

		select {
		case <-waiter.Wait():
		default:
			t.Fatalf("Flag not raised")
		}
	})

	t.Run("Second rise", func(t *testing.T) {
		// Second rise must be ignored
		assert.NoError(t, f.Raise(context.Background()))
		select {
		case <-waiter.Wait():
			t.Fatalf("Unexpected signal raised")
		default:
		}
	})

	t.Run("Reset", func(t *testing.T) {
		// Reset
		f.Reset()
		select {
		case <-waiter.Wait():
			t.Fatalf("Unexpected signal raised")
		default:
		}
		assert.False(t, f.IsRaised())
	})

	t.Run("Reraise", func(t *testing.T) {
		// Reraise
		assert.NoError(t, f.Raise(context.Background()))
		assert.True(t, f.IsRaised())

		select {
		case <-waiter.Wait():
		default:
			t.Fatalf("Flag not raised")
		}
	})

}

func TestFlagConcurrentSubscribe(t *testing.T) {
	concurrentWaiters := 10000

	f := NewFlag()
	gate := NewFlag()
	ch := make(chan *Waiter, concurrentWaiters)

	for i := 0; i < concurrentWaiters; i++ {
		gateWaiter := gate.Subscribe()
		go func() {
			assert.NoError(t, gateWaiter.WaitBlocking(context.Background()))
			ch <- f.Subscribe()
		}()
	}

	assert.NoError(t, gate.Raise(context.Background()))
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, f.Raise(context.Background()))

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
				assert.FailNow(t, "Waiter dont receive raise")
			}

		default:
			return
		}
	}
}
