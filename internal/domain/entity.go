package domain

import (
	"time"

	"github.com/google/uuid"
)

type ResourceType string

const (
	ResourceTypeFlightSeat ResourceType = "flight_seat"
	ResourceTypeHotelRoom  ResourceType = "hotel_room"
)

type InventoryStatus string

const (
	StatusAvailable InventoryStatus = "available"
	StatusHeld      InventoryStatus = "held"
	StatusBooked    InventoryStatus = "booked"
	StatusCancelled InventoryStatus = "cancelled"

	// InventoryStatus* aliases keep the domain API explicit at call sites.
	InventoryStatusAvailable InventoryStatus = StatusAvailable
	InventoryStatusHeld      InventoryStatus = StatusHeld
	InventoryStatusBooked    InventoryStatus = StatusBooked
	InventoryStatusCancelled InventoryStatus = StatusCancelled
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusExpired   BookingStatus = "expired"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Flight struct {
	ID             uuid.UUID `json:"id"`
	Code           string    `json:"code"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departure_time"`
	ArrivalTime    time.Time `json:"arrival_time"`
	BasePrice      float64   `json:"base_price"`
	AvailableSeats int       `json:"available_seats,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Hotel struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	City      string     `json:"city"`
	Address   string     `json:"address"`
	RoomTypes []RoomType `json:"room_types,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type RoomType struct {
	ID             uuid.UUID `json:"id"`
	HotelID        uuid.UUID `json:"hotel_id"`
	Name           string    `json:"name"`
	BasePrice      float64   `json:"base_price"`
	Capacity       int       `json:"capacity,omitempty"`
	AvailableRooms int       `json:"available_rooms,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type InventoryUnit struct {
	ID           uuid.UUID       `json:"id"`
	ResourceType ResourceType    `json:"resource_type"`
	ResourceID   uuid.UUID       `json:"resource_id"`
	UnitCode     string          `json:"unit_code"`
	Status       InventoryStatus `json:"status"`
	HeldUntil    *time.Time      `json:"held_until,omitempty"`
	Version      int             `json:"version"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type Booking struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"user_id"`
	Status         BookingStatus `json:"status"`
	IdempotencyKey string        `json:"idempotency_key,omitempty"`
	TotalAmount    float64       `json:"total_amount"`
	ExpiresAt      *time.Time    `json:"expires_at,omitempty"`
	ScheduledAt    *time.Time    `json:"scheduled_at,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Items          []BookingItem `json:"items,omitempty"`
}

type BookingItem struct {
	ID              uuid.UUID `json:"id"`
	BookingID       uuid.UUID `json:"booking_id"`
	InventoryUnitID uuid.UUID `json:"inventory_unit_id"`
	Price           float64   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
}

type IdempotencyKey struct {
	Key          string    `json:"key"`
	RequestHash  string    `json:"request_hash"`
	ResponseBody []byte    `json:"response_body"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Payment struct {
	ID          uuid.UUID     `json:"id"`
	BookingID   uuid.UUID     `json:"booking_id"`
	Status      PaymentStatus `json:"status"`
	Amount      float64       `json:"amount"`
	ProviderRef string        `json:"provider_ref,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type FlightSearchParams struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
	Date        string `form:"date" binding:"required"`
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
}

type HotelSearchParams struct {
	City     string `form:"city" binding:"required"`
	CheckIn  string `form:"checkin" binding:"required"`
	CheckOut string `form:"checkout" binding:"required"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type HoldSeatRequest struct {
	SeatID       uuid.UUID `json:"seat_id,omitempty"`
	HoldDuration int       `json:"hold_duration,omitempty" binding:"omitempty,min=1,max=60"` // minutes
}

type HoldRoomRequest struct {
	RoomID       uuid.UUID `json:"room_id,omitempty"`
	HoldDuration int       `json:"hold_duration,omitempty" binding:"omitempty,min=1,max=60"` // minutes
}

type CreateBookingRequest struct {
	Items       []CreateBookingItem `json:"items" binding:"required,min=1"`
	ScheduledAt *time.Time          `json:"scheduled_at,omitempty"`
}

type CreateBookingItem struct {
	InventoryUnitID uuid.UUID `json:"inventory_unit_id" binding:"required"`
}

type ConfirmBookingRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason"`
}

// ScheduleBookingRequest carries a new scheduled departure/check-in time for a
// booking. ScheduledAt must be in the future and the booking must belong to a
// terminal state (cancelled/expired) to be rejected.
type ScheduleBookingRequest struct {
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
