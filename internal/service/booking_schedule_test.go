package service

import (
	"errors"
	"testing"
	"time"

	"gobooking/internal/domain"
)

func TestScheduleBookingRejectsPastTime(t *testing.T) {
	service := &BookingService{}

	_, err := service.ScheduleBooking(nil, // ctx unused: validation short-circuits before any I/O
		nil,
		domain.ScheduleBookingRequest{ScheduledAt: time.Now().Add(-1 * time.Hour)},
	)
	if !errors.Is(err, ErrInvalidScheduleTime) {
		t.Fatalf("expected ErrInvalidScheduleTime, got %v", err)
	}
}

func TestScheduleBookingRejectsZeroTime(t *testing.T) {
	service := &BookingService{}

	_, err := service.ScheduleBooking(nil, nil, domain.ScheduleBookingRequest{})
	if !errors.Is(err, ErrInvalidScheduleTime) {
		t.Fatalf("expected ErrInvalidScheduleTime, got %v", err)
	}
}