package test

import (
	"signal"
	"testing"
)

func TestSignal(t *testing.T) {
	s := signal.NewSignal(1)
	waiter := s.Subscribe()
	s.Send()
	select {
	case <-waiter.Wait():
	default:
		t.Fatalf("Signal not received")
	}

	s.Send()
	s.Purge()

	select {
	case <-waiter.Wait():
		t.Fatalf("Unexpected signal received")
	default:
	}
	// Cancel waiter and try to send signal
	waiter.Cancel()
	s.Send()
	if !waiter.Empty() {
		t.Fatalf("Unexpected signal received")
	}
}
