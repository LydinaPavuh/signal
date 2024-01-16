package test

import (
	"signal"
	"testing"
)

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
