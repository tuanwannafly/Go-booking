package domain

import (
	"context"

	"github.com/google/uuid"
)

type SearchService interface {
	SearchFlights(ctx context.Context, params FlightSearchParams) ([]Flight, int64, error)
	SearchHotels(ctx context.Context, params HotelSearchParams) ([]Hotel, int64, error)
}

type HoldService interface {
	HoldSeat(ctx context.Context, seatID uuid.UUID, holdDurationMinutes int) (*InventoryUnit, error)
	HoldRoom(ctx context.Context, roomID uuid.UUID, holdDurationMinutes int) (*InventoryUnit, error)
	ReleaseHold(ctx context.Context, inventoryUnitID uuid.UUID) error
}

type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, req CreateBookingRequest, idempotencyKey string) (*Booking, error)
	GetBooking(ctx context.Context, id uuid.UUID) (*Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Booking, int64, error)
	ConfirmBooking(ctx context.Context, bookingID uuid.UUID, req ConfirmBookingRequest) (*Booking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID, req CancelBookingRequest) (*Booking, error)
	ScheduleBooking(ctx context.Context, bookingID uuid.UUID, req ScheduleBookingRequest) (*Booking, error)
}

type WorkerService interface {
	Start(ctx context.Context) error
	Stop() error
}

type RefundPolicy struct {
	FullRefundHours      int
	PartialRefundHours   int
	PartialRefundPercent float64
}

var DefaultRefundPolicy = RefundPolicy{
	FullRefundHours:      24,
	PartialRefundHours:   6,
	PartialRefundPercent: 0.5,
}
