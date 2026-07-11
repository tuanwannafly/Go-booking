package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/repository/postgres"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrBookingNotPending     = errors.New("booking is not in pending status")
	ErrBookingExpired        = errors.New("booking has expired")
	ErrIdempotencyConflict   = errors.New("idempotency key conflict: same key with different request")
	ErrIdempotentReplay      = errors.New("idempotent response replay")
	ErrInsufficientInventory = errors.New("insufficient inventory for booking")
	ErrPaymentFailed         = errors.New("payment failed")
)

type IdempotentReplayError struct {
	Booking *domain.Booking
}

func (e *IdempotentReplayError) Error() string { return ErrIdempotentReplay.Error() }

func (e *IdempotentReplayError) Unwrap() error { return ErrIdempotentReplay }

type BookingService struct {
	bookingRepo     *postgres.BookingRepository
	bookingItemRepo *postgres.BookingItemRepository
	inventoryRepo   *postgres.InventoryRepository
	idempotencyRepo *postgres.IdempotencyRepository
	paymentRepo     *postgres.PaymentRepository
	userRepo        *postgres.UserRepository
	holdService     *HoldService
	refundPolicy    domain.RefundPolicy
	bookingExpiry   time.Duration
}

func NewBookingService(
	bookingRepo *postgres.BookingRepository,
	bookingItemRepo *postgres.BookingItemRepository,
	inventoryRepo *postgres.InventoryRepository,
	idempotencyRepo *postgres.IdempotencyRepository,
	paymentRepo *postgres.PaymentRepository,
	userRepo *postgres.UserRepository,
	holdService *HoldService,
	refundPolicy domain.RefundPolicy,
	bookingExpiryMinutes int,
) *BookingService {
	return &BookingService{
		bookingRepo:     bookingRepo,
		bookingItemRepo: bookingItemRepo,
		inventoryRepo:   inventoryRepo,
		idempotencyRepo: idempotencyRepo,
		paymentRepo:     paymentRepo,
		userRepo:        userRepo,
		holdService:     holdService,
		refundPolicy:    refundPolicy,
		bookingExpiry:   time.Duration(bookingExpiryMinutes) * time.Minute,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID uuid.UUID, req domain.CreateBookingRequest, idempotencyKey string) (*domain.Booking, error) {
	// Check idempotency key if provided
	if idempotencyKey != "" {
		requestHash := s.computeRequestHash(req)
		existing, err := s.idempotencyRepo.Get(ctx, idempotencyKey)
		if err == nil {
			if existing.RequestHash != requestHash {
				return nil, ErrIdempotencyConflict
			}
			if existing.Status == "completed" {
				booking, decodeErr := decodeCachedBooking(existing.ResponseBody)
				if decodeErr != nil {
					return nil, decodeErr
				}
				return booking, &IdempotentReplayError{Booking: booking}
			}
			if existing.Status == "processing" {
				return s.waitForIdempotencyCompletion(ctx, idempotencyKey, requestHash)
			}
		}
		if err != nil && !errors.Is(err, postgres.ErrIdempotencyKeyNotFound) {
			return nil, err
		}

		// Create processing entry
		if err := s.idempotencyRepo.Create(ctx, &domain.IdempotencyKey{
			Key:          idempotencyKey,
			RequestHash:  requestHash,
			ResponseBody: nil,
			Status:       "processing",
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}); err != nil {
			if errors.Is(err, postgres.ErrIdempotencyKeyConflict) {
				return s.waitForIdempotencyCompletion(ctx, idempotencyKey, requestHash)
			}
			return nil, err
		}
	}

	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Validate inventory units are held by this user
	// In a real implementation, you'd check the hold ownership
	inventoryUnitIDs := make([]uuid.UUID, len(req.Items))
	for i, item := range req.Items {
		inventoryUnitIDs[i] = item.InventoryUnitID
	}

	// Check all units are held
	units, err := s.inventoryRepo.GetByIDs(ctx, inventoryUnitIDs)
	if err != nil {
		return nil, err
	}

	if len(units) != len(inventoryUnitIDs) {
		return nil, ErrInsufficientInventory
	}

	// Calculate total amount
	var totalAmount float64
	for _, unit := range units {
		if unit.Status != domain.InventoryStatusHeld {
			return nil, ErrInsufficientInventory
		}
		// In real implementation, price would come from the resource (flight/hotel)
		totalAmount += 100.0 // placeholder
	}

	// Create booking
	booking := &domain.Booking{
		ID:             uuid.New(),
		UserID:         userID,
		Status:         domain.BookingStatusPending,
		IdempotencyKey: idempotencyKey,
		TotalAmount:    totalAmount,
		ExpiresAt:      ptr(time.Now().Add(s.bookingExpiry)),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// Create booking items
	for _, item := range req.Items {
		bookingItem := &domain.BookingItem{
			ID:              uuid.New(),
			BookingID:       booking.ID,
			InventoryUnitID: item.InventoryUnitID,
			Price:           100.0, // placeholder
			CreatedAt:       time.Now(),
		}
		if err := s.bookingItemRepo.Create(ctx, bookingItem); err != nil {
			return nil, err
		}
		booking.Items = append(booking.Items, *bookingItem)
	}

	// Save idempotency response
	if idempotencyKey != "" {
		if err := s.idempotencyRepo.Update(ctx, &domain.IdempotencyKey{
			Key:          idempotencyKey,
			RequestHash:  s.computeRequestHash(req),
			ResponseBody: mustJSON(booking),
			Status:       "completed",
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}); err != nil {
			return nil, err
		}
	}

	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return s.bookingRepo.GetByIDWithItems(ctx, id)
}

func (s *BookingService) GetBookingByIdempotencyKey(ctx context.Context, key string) (*domain.Booking, error) {
	return s.bookingRepo.GetByIdempotencyKey(ctx, key)
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]domain.Booking, int64, error) {
	return s.bookingRepo.GetByUserID(ctx, userID, page, pageSize)
}

func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID uuid.UUID, req domain.ConfirmBookingRequest) (*domain.Booking, error) {
	if req.PaymentMethod == "" {
		return nil, errors.New("payment method is required")
	}
	booking, err := s.bookingRepo.GetByIDWithItems(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	if booking.Status != domain.BookingStatusPending {
		if booking.Status == domain.BookingStatusConfirmed {
			return booking, &IdempotentReplayError{Booking: booking}
		}
		return nil, ErrBookingNotPending
	}

	// Check if expired
	if booking.ExpiresAt != nil && time.Now().After(*booking.ExpiresAt) {
		// Expire the booking
		s.expireBooking(ctx, booking)
		return nil, ErrBookingExpired
	}

	// Create mock payment
	payment := &domain.Payment{
		ID:          uuid.New(),
		BookingID:   bookingID,
		Status:      domain.PaymentStatusCompleted,
		Amount:      booking.TotalAmount,
		ProviderRef: fmt.Sprintf("mock_%s", uuid.New().String()[:8]),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.bookingRepo.Confirm(ctx, bookingID, payment); err != nil {
		return nil, err
	}

	booking.Status = domain.BookingStatusConfirmed
	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID uuid.UUID, req domain.CancelBookingRequest) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByIDWithItems(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	if booking.Status == domain.BookingStatusCancelled {
		return nil, errors.New("booking already cancelled")
	}

	if booking.Status == domain.BookingStatusExpired {
		return nil, errors.New("booking already expired")
	}

	// Calculate refund based on policy
	refundAmount := s.calculateRefund(booking)

	// Get inventory unit IDs
	inventoryUnitIDs := make([]uuid.UUID, len(booking.Items))
	for i, item := range booking.Items {
		inventoryUnitIDs[i] = item.InventoryUnitID
	}

	// Release inventory based on current status
	if booking.Status == domain.BookingStatusConfirmed {
		// Booked -> Cancelled (refund)
		if err := s.inventoryRepo.CancelBooking(ctx, inventoryUnitIDs); err != nil {
			return nil, err
		}
		// Update payment status to refunded
		if err := s.paymentRepo.UpdateStatusByBookingID(ctx, bookingID, domain.PaymentStatusRefunded, fmt.Sprintf("refund_%.2f", refundAmount)); err != nil {
			return nil, err
		}
	} else if booking.Status == domain.BookingStatusPending {
		// Pending -> Expired (release holds)
		if err := s.inventoryRepo.ExpireBooking(ctx, inventoryUnitIDs); err != nil {
			return nil, err
		}
	}

	// Update booking status
	newStatus := domain.BookingStatusCancelled
	if booking.Status == domain.BookingStatusPending {
		newStatus = domain.BookingStatusExpired
	}
	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, newStatus); err != nil {
		return nil, err
	}

	booking.Status = newStatus
	return booking, nil
}

func (s *BookingService) ExpirePendingBookings(ctx context.Context) (int, error) {
	bookings, err := s.bookingRepo.GetExpiredPendingBookings(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, booking := range bookings {
		inventoryUnitIDs := make([]uuid.UUID, len(booking.Items))
		for i, item := range booking.Items {
			inventoryUnitIDs[i] = item.InventoryUnitID
		}

		if err := s.inventoryRepo.ExpireBooking(ctx, inventoryUnitIDs); err != nil {
			continue
		}

		if err := s.bookingRepo.UpdateStatus(ctx, booking.ID, domain.BookingStatusExpired); err != nil {
			continue
		}
		count++
	}

	return count, nil
}

func (s *BookingService) expireBooking(ctx context.Context, booking *domain.Booking) error {
	inventoryUnitIDs := make([]uuid.UUID, len(booking.Items))
	for i, item := range booking.Items {
		inventoryUnitIDs[i] = item.InventoryUnitID
	}

	if err := s.inventoryRepo.ExpireBooking(ctx, inventoryUnitIDs); err != nil {
		return err
	}

	return s.bookingRepo.UpdateStatus(ctx, booking.ID, domain.BookingStatusExpired)
}

func (s *BookingService) calculateRefund(booking *domain.Booking) float64 {
	hoursUntilDeparture := 24.0 // placeholder - in real implementation, calculate from flight/hotel time

	if hoursUntilDeparture > float64(s.refundPolicy.FullRefundHours) {
		return booking.TotalAmount
	} else if hoursUntilDeparture > float64(s.refundPolicy.PartialRefundHours) {
		return booking.TotalAmount * s.refundPolicy.PartialRefundPercent
	}
	return 0
}

func (s *BookingService) computeRequestHash(req domain.CreateBookingRequest) string {
	data := mustJSON(req)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *BookingService) waitForIdempotencyCompletion(ctx context.Context, key, requestHash string) (*domain.Booking, error) {
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(2 * time.Second)
	defer deadline.Stop()

	for {
		select {
		case <-ticker.C:
			record, err := s.idempotencyRepo.Get(ctx, key)
			if err != nil {
				if errors.Is(err, postgres.ErrIdempotencyKeyNotFound) {
					continue
				}
				return nil, err
			}
			if record.RequestHash != requestHash {
				return nil, ErrIdempotencyConflict
			}
			if record.Status == "completed" {
				booking, err := decodeCachedBooking(record.ResponseBody)
				if err != nil {
					return nil, err
				}
				return booking, &IdempotentReplayError{Booking: booking}
			}
		case <-deadline.C:
			return nil, fmt.Errorf("idempotency request still processing")
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func decodeCachedBooking(data []byte) (*domain.Booking, error) {
	var booking domain.Booking
	if err := json.Unmarshal(data, &booking); err != nil {
		return nil, fmt.Errorf("decode cached booking: %w", err)
	}
	return &booking, nil
}

func ptr[T any](v T) *T {
	return &v
}

func mustJSON(value interface{}) []byte {
	data, _ := json.Marshal(value)
	return data
}
