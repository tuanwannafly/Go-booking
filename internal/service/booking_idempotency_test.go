package service

import (
	"encoding/json"
	"errors"
	"testing"

	"gobooking/internal/domain"
)

func TestDecodeCachedBooking(t *testing.T) {
	want := &domain.Booking{Status: domain.BookingStatusPending, TotalAmount: 100}
	payload, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := decodeCachedBooking(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != want.Status || got.TotalAmount != want.TotalAmount {
		t.Fatalf("decoded booking mismatch: got %#v want %#v", got, want)
	}
}

func TestDecodeCachedBookingRejectsInvalidPayload(t *testing.T) {
	_, err := decodeCachedBooking([]byte("not-json"))
	if err == nil {
		t.Fatal("expected invalid payload error")
	}
}

func TestIdempotentReplayErrorUnwraps(t *testing.T) {
	err := &IdempotentReplayError{Booking: &domain.Booking{}}
	if !errors.Is(err, ErrIdempotentReplay) {
		t.Fatal("expected replay error to unwrap to ErrIdempotentReplay")
	}
}
