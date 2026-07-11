package concurrency

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"gobooking/internal/domain"
	"gobooking/internal/repository/postgres"
	redisrepo "gobooking/internal/repository/redis"
	"gobooking/internal/service"
)

func TestHoldSeat_ConcurrentRequests(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	// Setup test database and Redis
	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	// Create test flight and seats
	flightID := createTestFlight(t, db)
	seatIDs := createTestSeats(t, db, flightID, 1) // Only 1 seat for contention test

	// Create services
	inventoryRepo := postgres.NewInventoryRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)

	// Test: 100 goroutines trying to hold the same seat
	const numGoroutines = 100
	var successCount int32
	var failureCount int32
	var wg sync.WaitGroup

	ctx := context.Background()
	seatID := seatIDs[0]

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			unit, err := holdService.HoldSeat(ctx, seatID, 10)
			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				return
			}

			if unit != nil && unit.Status == domain.InventoryStatusHeld {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failureCount, 1)
			}
		}()
	}

	wg.Wait()

	t.Logf("Results: %d succeeded, %d failed", successCount, failureCount)

	// Assert: only 1 should succeed
	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}

	if failureCount != numGoroutines-1 {
		t.Errorf("Expected %d failures, got %d", numGoroutines-1, failureCount)
	}

	// Verify the seat is actually held
	unit, err := inventoryRepo.GetByID(ctx, seatID)
	if err != nil {
		t.Fatalf("Failed to get seat: %v", err)
	}
	if unit.Status != domain.InventoryStatusHeld {
		t.Errorf("Expected seat status 'held', got '%s'", unit.Status)
	}
}

func TestHoldRoom_ConcurrentRequests(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	hotelID := createTestHotel(t, db)
	roomTypeID := createTestRoomType(t, db, hotelID)
	roomIDs := createTestRooms(t, db, roomTypeID, 1) // Only 1 room

	inventoryRepo := postgres.NewInventoryRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)

	const numGoroutines = 50
	var successCount int32
	var failureCount int32
	var wg sync.WaitGroup

	ctx := context.Background()
	roomID := roomIDs[0]

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			unit, err := holdService.HoldRoom(ctx, roomID, 10)
			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				return
			}

			if unit != nil && unit.Status == domain.InventoryStatusHeld {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failureCount, 1)
			}
		}()
	}

	wg.Wait()

	t.Logf("Results: %d succeeded, %d failed", successCount, failureCount)

	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}

	if failureCount != numGoroutines-1 {
		t.Errorf("Expected %d failures, got %d", numGoroutines-1, failureCount)
	}
}

func TestBooking_ConcurrentCreation(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	// Create user
	userID := createTestUser(t, db)

	// Create flight with 1 seat
	flightID := createTestFlight(t, db)
	seatIDs := createTestSeats(t, db, flightID, 1)
	seatID := seatIDs[0]

	// Hold the seat first
	inventoryRepo := postgres.NewInventoryRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)

	ctx := context.Background()
	heldUnit, err := holdService.HoldSeat(ctx, seatID, 10)
	if err != nil {
		t.Fatalf("Failed to hold seat: %v", err)
	}

	// Create services
	bookingRepo := postgres.NewBookingRepository(db)
	bookingItemRepo := postgres.NewBookingItemRepository(db)
	idempotencyRepo := postgres.NewIdempotencyRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	userRepo := postgres.NewUserRepository(db)

	bookingService := service.NewBookingService(
		bookingRepo,
		bookingItemRepo,
		inventoryRepo,
		idempotencyRepo,
		paymentRepo,
		userRepo,
		holdService,
		domain.DefaultRefundPolicy,
		10,
	)

	// Test: 10 goroutines trying to create booking with same held seat
	const numGoroutines = 10
	var successCount int32
	var failureCount int32
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := domain.CreateBookingRequest{
				Items: []domain.CreateBookingItem{
					{InventoryUnitID: heldUnit.ID},
				},
			}

			idempotencyKey := uuid.New().String()
			_, err := bookingService.CreateBooking(ctx, userID, req, idempotencyKey)
			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}()
	}

	wg.Wait()

	t.Logf("Results: %d succeeded, %d failed", successCount, failureCount)

	// Only 1 should succeed
	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}
}

func TestIdempotency_ConcurrentSameKey(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userID := createTestUser(t, db)
	flightID := createTestFlight(t, db)
	seatIDs := createTestSeats(t, db, flightID, 1)
	seatID := seatIDs[0]

	inventoryRepo := postgres.NewInventoryRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)

	ctx := context.Background()
	heldUnit, err := holdService.HoldSeat(ctx, seatID, 10)
	if err != nil {
		t.Fatalf("Failed to hold seat: %v", err)
	}

	bookingRepo := postgres.NewBookingRepository(db)
	bookingItemRepo := postgres.NewBookingItemRepository(db)
	idempotencyRepo := postgres.NewIdempotencyRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	userRepo := postgres.NewUserRepository(db)

	bookingService := service.NewBookingService(
		bookingRepo,
		bookingItemRepo,
		inventoryRepo,
		idempotencyRepo,
		paymentRepo,
		userRepo,
		holdService,
		domain.DefaultRefundPolicy,
		10,
	)

	// Test: 10 goroutines with SAME idempotency key
	const numGoroutines = 10
	var successCount int32
	var conflictCount int32
	var otherErrorCount int32
	var wg sync.WaitGroup

	idempotencyKey := "test-idempotency-key-123"

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := domain.CreateBookingRequest{
				Items: []domain.CreateBookingItem{
					{InventoryUnitID: heldUnit.ID},
				},
			}

			_, err := bookingService.CreateBooking(ctx, userID, req, idempotencyKey)
			if err != nil {
				errMsg := err.Error()
				if errMsg == "idempotency key conflict: different request body" {
					atomic.AddInt32(&conflictCount, 1)
				} else if errMsg == "idempotent response cached" {
					atomic.AddInt32(&successCount, 1)
				} else {
					atomic.AddInt32(&otherErrorCount, 1)
				}
				return
			}

			atomic.AddInt32(&successCount, 1)
		}()
	}

	wg.Wait()

	t.Logf("Results: %d succeeded, %d conflicts, %d other errors", successCount, conflictCount, otherErrorCount)

	// Should have exactly 1 success (first request) and rest should get cached response or conflict
	if successCount != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount)
	}
}

// Helper functions

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://gobooking:gobooking@localhost:5432/gobooking?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

func createTestUser(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()
	userID := uuid.New()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)
	`, userID, "test@example.com", "hash", "user")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	return userID
}

func createTestFlight(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()
	flightID := uuid.New()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO flights (id, code, origin, destination, departure_time, arrival_time, base_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, flightID, "TEST123", "SGN", "HAN", time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour), 1000000)
	if err != nil {
		t.Fatalf("Failed to create flight: %v", err)
	}
	return flightID
}

func createTestSeats(t *testing.T, db *pgxpool.Pool, flightID uuid.UUID, count int) []uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var seatIDs []uuid.UUID

	for i := 0; i < count; i++ {
		seatID := uuid.New()
		seatIDs = append(seatIDs, seatID)
		unitCode := fmt.Sprintf("Y%dA", i+1)

		_, err := db.Exec(ctx, `
			INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
			VALUES ($1, 'flight_seat', $2, $3, 'available', 1)
		`, seatID, flightID, unitCode)

		if err != nil {
			t.Fatalf("Failed to create seat: %v", err)
		}
	}

	return seatIDs
}

func createTestHotel(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()
	hotelID := uuid.New()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO hotels (id, name, city, address)
		VALUES ($1, $2, $3, $4)
	`, hotelID, "Test Hotel", "Test City", "Test Address")
	if err != nil {
		t.Fatalf("Failed to create hotel: %v", err)
	}
	return hotelID
}

func createTestRoomType(t *testing.T, db *pgxpool.Pool, hotelID uuid.UUID) uuid.UUID {
	t.Helper()
	rtID := uuid.New()
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO room_types (id, hotel_id, name, base_price, capacity)
		VALUES ($1, $2, $3, $4, $5)
	`, rtID, hotelID, "Standard", 1000000, 2)
	if err != nil {
		t.Fatalf("Failed to create room type: %v", err)
	}
	return rtID
}

func createTestRooms(t *testing.T, db *pgxpool.Pool, roomTypeID uuid.UUID, count int) []uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var roomIDs []uuid.UUID

	for i := 0; i < count; i++ {
		roomID := uuid.New()
		roomIDs = append(roomIDs, roomID)
		unitCode := fmt.Sprintf("10%d", i+1)

		_, err := db.Exec(ctx, `
			INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
			VALUES ($1, 'hotel_room', $2, $3, 'available', 1)
		`, roomID, roomTypeID, unitCode)

		if err != nil {
			t.Fatalf("Failed to create room: %v", err)
		}
	}

	return roomIDs
}
