package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type FlightRepository interface {
	Create(ctx context.Context, flight *Flight) error
	GetByID(ctx context.Context, id uuid.UUID) (*Flight, error)
	Search(ctx context.Context, params FlightSearchParams) ([]Flight, int64, error)
}

type HotelRepository interface {
	Create(ctx context.Context, hotel *Hotel) error
	GetByID(ctx context.Context, id uuid.UUID) (*Hotel, error)
	Search(ctx context.Context, params HotelSearchParams) ([]Hotel, int64, error)
}

type RoomTypeRepository interface {
	Create(ctx context.Context, rt *RoomType) error
	GetByID(ctx context.Context, id uuid.UUID) (*RoomType, error)
	GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]RoomType, error)
}

type InventoryRepository interface {
	Create(ctx context.Context, unit *InventoryUnit) error
	GetByID(ctx context.Context, id uuid.UUID) (*InventoryUnit, error)
	GetByResource(ctx context.Context, resourceType ResourceType, resourceID uuid.UUID) ([]InventoryUnit, error)
	GetAvailableByResource(ctx context.Context, resourceType ResourceType, resourceID uuid.UUID) ([]InventoryUnit, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status InventoryStatus, heldUntil *time.Time, version int) (bool, error)
	UpdateStatusOptimistic(ctx context.Context, id uuid.UUID, status InventoryStatus, heldUntil *time.Time, version int) (bool, error)
	ReleaseExpiredHolds(ctx context.Context) (int64, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Booking, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status BookingStatus) error
	UpdateExpiresAt(ctx context.Context, id uuid.UUID, expiresAt *time.Time) error
	GetByIdempotencyKey(ctx context.Context, key string) (*Booking, error)
}

type BookingItemRepository interface {
	Create(ctx context.Context, item *BookingItem) error
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]BookingItem, error)
}

type IdempotencyRepository interface {
	Create(ctx context.Context, key *IdempotencyKey) error
	Get(ctx context.Context, key string) (*IdempotencyKey, error)
	Update(ctx context.Context, key *IdempotencyKey) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus) error
}
