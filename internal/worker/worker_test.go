package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHoldReleaseWorkerRejectsInvalidInterval(t *testing.T) {
	worker := NewHoldReleaseWorker(nil, 0)
	if err := worker.Start(context.Background()); !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
	if err := worker.Stop(); err != nil {
		t.Fatalf("stop returned unexpected error: %v", err)
	}
	if err := worker.Stop(); err != nil {
		t.Fatalf("second stop returned unexpected error: %v", err)
	}
}

func TestBookingExpireWorkerRejectsInvalidInterval(t *testing.T) {
	worker := NewBookingExpireWorker(nil, -time.Second)
	if err := worker.Start(context.Background()); !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestWorkerServiceStopIsIdempotent(t *testing.T) {
	service := NewWorkerService()
	if err := service.Start(context.Background()); err != nil {
		t.Fatalf("start returned unexpected error: %v", err)
	}
	if err := service.Stop(); err != nil {
		t.Fatalf("stop returned unexpected error: %v", err)
	}
	if err := service.Stop(); err != nil {
		t.Fatalf("second stop returned unexpected error: %v", err)
	}
}
