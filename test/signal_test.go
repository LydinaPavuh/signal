package test

import (
	"signal/src/signal"
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

func TestFlag(t *testing.T) {
	f := signal.NewFlag()
	waiter := f.Subscribe()
	f.Raise()
	if !f.IsRaised() {
		t.Fatalf("Flag not raised")
	}
	select {
	case <-waiter.Wait():
	default:
		t.Fatalf("Flag not raised")
	}

	// Second rise must be ignored
	f.Raise()
	select {
	case <-waiter.Wait():
		t.Fatalf("Unexpected signal raised")
	default:
	}

	// Reset
	f.Reset()
	select {
	case <-waiter.Wait():
		t.Fatalf("Unexpected signal raised")
	default:
	}
	if f.IsRaised() {
		t.Fatalf("Flag do not be reseted")
	}

	// Reraise
	f.Raise()
	if !f.IsRaised() {
		t.Fatalf("Flag not raised")
	}
	select {
	case <-waiter.Wait():
	default:
		t.Fatalf("Flag not raised")
	}
}
