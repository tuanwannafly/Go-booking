package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"gobooking/internal/domain"
)

func TestConfirmBookingRequestRequiresPaymentMethod(t *testing.T) {
	service := &BookingService{}
	_, err := service.ConfirmBooking(context.Background(), uuid.Nil, domain.ConfirmBookingRequest{})
	if err == nil || err.Error() != "payment method is required" {
		t.Fatalf("expected payment method validation error, got %v", err)
	}
}
