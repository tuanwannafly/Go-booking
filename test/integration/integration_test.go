package integration

import (
	"context"
	"fmt"
	"os"
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

func TestFullFlow_SearchHoldBookConfirmCancel(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	// Setup services
	flightRepo := postgres.NewFlightRepository(db)
	hotelRepo := postgres.NewHotelRepository(db)
	roomTypeRepo := postgres.NewRoomTypeRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	bookingItemRepo := postgres.NewBookingItemRepository(db)
	idempotencyRepo := postgres.NewIdempotencyRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)

	searchService := service.NewSearchService(flightRepo, hotelRepo, roomTypeRepo, inventoryRepo)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)
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

	ctx := context.Background()

	// 1. Create test user
	userID := createTestUser(t, db)

	// 2. Create test flight with seats
	flightID := createTestFlight(t, db)
	seatIDs := createTestSeats(t, db, flightID, 5)
	seatID := seatIDs[0]

	// 3. Search flights
	_, total, err := searchService.SearchFlights(ctx, domain.FlightSearchParams{
		Origin:      "SGN",
		Destination: "HAN",
		Date:        time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		t.Fatalf("Search flights failed: %v", err)
	}
	if total == 0 {
		t.Fatal("Expected at least 1 flight")
	}
	t.Logf("Found %d flights", total)

	// 4. Hold a seat (pessimistic lock)
	heldUnit, err := holdService.HoldSeat(ctx, seatID, 10)
	if err != nil {
		t.Fatalf("Hold seat failed: %v", err)
	}
	if heldUnit.Status != domain.InventoryStatusHeld {
		t.Fatalf("Expected held status, got %s", heldUnit.Status)
	}
	t.Logf("Seat held: %s", heldUnit.ID)

	// 5. Create booking with idempotency key
	idempotencyKey := "test-idempotency-" + uuid.New().String()
	booking, err := bookingService.CreateBooking(ctx, userID, domain.CreateBookingRequest{
		Items: []domain.CreateBookingItem{
			{InventoryUnitID: heldUnit.ID},
		},
	}, idempotencyKey)
	if err != nil {
		t.Fatalf("Create booking failed: %v", err)
	}
	if booking.Status != domain.BookingStatusPending {
		t.Fatalf("Expected pending status, got %s", booking.Status)
	}
	t.Logf("Booking created: %s", booking.ID)

	// 6. Confirm booking (mock payment)
	confirmedBooking, err := bookingService.ConfirmBooking(ctx, booking.ID, domain.ConfirmBookingRequest{
		PaymentMethod: "mock",
	})
	if err != nil {
		t.Fatalf("Confirm booking failed: %v", err)
	}
	if confirmedBooking.Status != domain.BookingStatusConfirmed {
		t.Fatalf("Expected confirmed status, got %s", confirmedBooking.Status)
	}
	t.Logf("Booking confirmed: %s", confirmedBooking.ID)

	// 7. Verify inventory is now booked
	unit, err := inventoryRepo.GetByID(ctx, heldUnit.ID)
	if err != nil {
		t.Fatalf("Get inventory failed: %v", err)
	}
	if unit.Status != domain.InventoryStatusBooked {
		t.Fatalf("Expected booked status, got %s", unit.Status)
	}

	// 8. Cancel booking (should refund)
	cancelledBooking, err := bookingService.CancelBooking(ctx, booking.ID, domain.CancelBookingRequest{
		Reason: "Changed plans",
	})
	if err != nil {
		t.Fatalf("Cancel booking failed: %v", err)
	}
	if cancelledBooking.Status != domain.BookingStatusCancelled {
		t.Fatalf("Expected cancelled status, got %s", cancelledBooking.Status)
	}
	t.Logf("Booking cancelled: %s", cancelledBooking.ID)

	// 9. Verify inventory is available again
	unit, err = inventoryRepo.GetByID(ctx, heldUnit.ID)
	if err != nil {
		t.Fatalf("Get inventory failed: %v", err)
	}
	if unit.Status != domain.InventoryStatusAvailable {
		t.Fatalf("Expected available status after cancel, got %s", unit.Status)
	}

	t.Log("Full flow test passed!")
}

func TestHotelFlow_SearchHoldBook(t *testing.T) {
	if testing.Short() || os.Getenv("GBOOKING_INTEGRATION") != "1" {
		t.Skip("requires PostgreSQL and Redis; set GBOOKING_INTEGRATION=1")
	}

	db := setupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	flightRepo := postgres.NewFlightRepository(db)
	hotelRepo := postgres.NewHotelRepository(db)
	roomTypeRepo := postgres.NewRoomTypeRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)
	bookingItemRepo := postgres.NewBookingItemRepository(db)
	idempotencyRepo := postgres.NewIdempotencyRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	userRepo := postgres.NewUserRepository(db)
	redisRepo := redisrepo.NewRedisClient("localhost:6379", "", 0)

	searchService := service.NewSearchService(flightRepo, hotelRepo, roomTypeRepo, inventoryRepo)
	holdService := service.NewHoldService(inventoryRepo, redisRepo, 10)
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

	ctx := context.Background()

	// Create test user
	userID := createTestUser(t, db)

	// Create test hotel with rooms
	hotelID := createTestHotel(t, db)
	roomTypeID := createTestRoomType(t, db, hotelID)
	roomIDs := createTestRooms(t, db, roomTypeID, 5)
	roomID := roomIDs[0]

	// Search hotels
	_, total, err := searchService.SearchHotels(ctx, domain.HotelSearchParams{
		City:     "Test City",
		CheckIn:  time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		CheckOut: time.Now().Add(48 * time.Hour).Format("2006-01-02"),
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("Search hotels failed: %v", err)
	}
	if total == 0 {
		t.Fatal("Expected at least 1 hotel")
	}
	t.Logf("Found %d hotels", total)

	// Hold a room (optimistic lock)
	heldUnit, err := holdService.HoldRoom(ctx, roomID, 10)
	if err != nil {
		t.Fatalf("Hold room failed: %v", err)
	}
	if heldUnit.Status != domain.InventoryStatusHeld {
		t.Fatalf("Expected held status, got %s", heldUnit.Status)
	}
	t.Logf("Room held: %s", heldUnit.ID)

	// Create booking
	booking, err := bookingService.CreateBooking(ctx, userID, domain.CreateBookingRequest{
		Items: []domain.CreateBookingItem{
			{InventoryUnitID: heldUnit.ID},
		},
	}, "hotel-idempotency-"+uuid.New().String())
	if err != nil {
		t.Fatalf("Create booking failed: %v", err)
	}
	t.Logf("Hotel booking created: %s", booking.ID)

	// Confirm booking
	confirmedBooking, err := bookingService.ConfirmBooking(ctx, booking.ID, domain.ConfirmBookingRequest{
		PaymentMethod: "mock",
	})
	if err != nil {
		t.Fatalf("Confirm booking failed: %v", err)
	}
	if confirmedBooking.Status != domain.BookingStatusConfirmed {
		t.Fatalf("Expected confirmed status, got %s", confirmedBooking.Status)
	}

	// Verify inventory is booked
	unit, err := inventoryRepo.GetByID(ctx, heldUnit.ID)
	if err != nil {
		t.Fatalf("Get inventory failed: %v", err)
	}
	if unit.Status != domain.InventoryStatusBooked {
		t.Fatalf("Expected booked status, got %s", unit.Status)
	}

	t.Log("Hotel flow test passed!")
}

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
