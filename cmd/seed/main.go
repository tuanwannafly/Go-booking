package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	db, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Seed flights
	flightIDs, err := seedFlights(ctx, db)
	if err != nil {
		log.Fatalf("Failed to seed flights: %v", err)
	}
	log.Printf("Created %d flights", len(flightIDs))

	// Seed hotels
	hotelIDs, err := seedHotels(ctx, db)
	if err != nil {
		log.Fatalf("Failed to seed hotels: %v", err)
	}
	log.Printf("Created %d hotels", len(hotelIDs))

	// Seed room types and rooms
	for _, hotelID := range hotelIDs {
		roomTypeIDs, err := seedRoomTypes(ctx, db, hotelID)
		if err != nil {
			log.Fatalf("Failed to seed room types: %v", err)
		}
		for _, rtID := range roomTypeIDs {
			if err := seedRooms(ctx, db, rtID, 10); err != nil {
				log.Fatalf("Failed to seed rooms: %v", err)
			}
		}
	}

	// Seed seats for flights
	for _, flightID := range flightIDs {
		if err := seedSeats(ctx, db, flightID, 180); err != nil {
			log.Fatalf("Failed to seed seats: %v", err)
		}
	}

	log.Println("Seeding completed successfully!")
}

func seedFlights(ctx context.Context, db *pgxpool.Pool) ([]uuid.UUID, error) {
	routes := []struct {
		code        string
		origin      string
		destination string
		basePrice   float64
		duration    time.Duration
	}{
		{"VN100", "SGN", "HAN", 1200000, 2 * time.Hour},
		{"VN101", "HAN", "SGN", 1200000, 2 * time.Hour},
		{"VN200", "SGN", "DAD", 900000, 1*time.Hour + 30*time.Minute},
		{"VN201", "DAD", "SGN", 900000, 1*time.Hour + 30*time.Minute},
		{"VN300", "HAN", "DAD", 800000, 1*time.Hour + 15*time.Minute},
		{"VN301", "DAD", "HAN", 800000, 1*time.Hour + 15*time.Minute},
		{"VN400", "SGN", "PQC", 1500000, 1 * time.Hour},
		{"VN401", "PQC", "SGN", 1500000, 1 * time.Hour},
		{"VN500", "HAN", "PQC", 1800000, 2 * time.Hour},
		{"VN501", "PQC", "HAN", 1800000, 2 * time.Hour},
	}

	var flightIDs []uuid.UUID
	baseTime := time.Now().Add(24 * time.Hour).Truncate(time.Hour)

	for i, route := range routes {
		flightID := uuid.New()
		flightIDs = append(flightIDs, flightID)

		departureTime := baseTime.Add(time.Duration(i*2) * time.Hour)
		arrivalTime := departureTime.Add(route.duration)

		_, err := db.Exec(ctx, `
			INSERT INTO flights (id, code, origin, destination, departure_time, arrival_time, base_price)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, flightID, route.code, route.origin, route.destination, departureTime, arrivalTime, route.basePrice)

		if err != nil {
			return nil, fmt.Errorf("insert flight %s: %w", route.code, err)
		}
	}

	return flightIDs, nil
}

func seedHotels(ctx context.Context, db *pgxpool.Pool) ([]uuid.UUID, error) {
	hotels := []struct {
		name    string
		city    string
		address string
	}{
		{"Grand Saigon Hotel", "Ho Chi Minh City", "123 Dong Khoi, District 1"},
		{"Hanoi Plaza Hotel", "Hanoi", "456 Le Loi, Hoan Kiem"},
		{"Da Nang Beach Resort", "Da Nang", "789 Vo Nguyen Giap, Son Tra"},
		{"Phu Quoc Paradise", "Phu Quoc", "321 Bai Truong, Duong Dong"},
		{"Nha Trang Ocean View", "Nha Trang", "555 Tran Phu, Loc Tho"},
	}

	var hotelIDs []uuid.UUID
	for _, h := range hotels {
		hotelID := uuid.New()
		hotelIDs = append(hotelIDs, hotelID)

		_, err := db.Exec(ctx, `
			INSERT INTO hotels (id, name, city, address)
			VALUES ($1, $2, $3, $4)
		`, hotelID, h.name, h.city, h.address)

		if err != nil {
			return nil, fmt.Errorf("insert hotel %s: %w", h.name, err)
		}
	}

	return hotelIDs, nil
}

func seedRoomTypes(ctx context.Context, db *pgxpool.Pool, hotelID uuid.UUID) ([]uuid.UUID, error) {
	roomTypes := []struct {
		name      string
		basePrice float64
		capacity  int
	}{
		{"Standard", 800000, 2},
		{"Deluxe", 1200000, 2},
		{"Suite", 2500000, 4},
		{"Family", 1800000, 4},
	}

	var rtIDs []uuid.UUID
	for _, rt := range roomTypes {
		rtID := uuid.New()
		rtIDs = append(rtIDs, rtID)

		_, err := db.Exec(ctx, `
			INSERT INTO room_types (id, hotel_id, name, base_price, capacity)
			VALUES ($1, $2, $3, $4, $5)
		`, rtID, hotelID, rt.name, rt.basePrice, rt.capacity)

		if err != nil {
			return nil, fmt.Errorf("insert room type %s: %w", rt.name, err)
		}
	}

	return rtIDs, nil
}

func seedRooms(ctx context.Context, db *pgxpool.Pool, roomTypeID uuid.UUID, count int) error {
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	floors := []int{1, 2, 3, 4, 5}

	for i := 0; i < count; i++ {
		roomID := uuid.New()
		floor := floors[rand.Intn(len(floors))]
		row := rows[rand.Intn(len(rows))]
		unitCode := fmt.Sprintf("%d%s%02d", floor, row, (i%10)+1)

		_, err := db.Exec(ctx, `
			INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
			VALUES ($1, 'hotel_room', $2, $3, 'available', 1)
		`, roomID, roomTypeID, unitCode)

		if err != nil {
			return fmt.Errorf("insert room %s: %w", unitCode, err)
		}
	}
	return nil
}

func seedSeats(ctx context.Context, db *pgxpool.Pool, flightID uuid.UUID, count int) error {
	rows := []string{"A", "B", "C", "D", "E", "F"}
	classes := []struct {
		prefix string
		count  int
	}{
		{"J", 12},         // Business class
		{"Y", count - 12}, // Economy
	}

	seatNum := 1
	for _, class := range classes {
		for i := 0; i < class.count; i++ {
			seatID := uuid.New()
			row := rows[seatNum%len(rows)]
			unitCode := fmt.Sprintf("%s%d%s", class.prefix, (seatNum/len(rows))+1, row)

			_, err := db.Exec(ctx, `
				INSERT INTO inventory_units (id, resource_type, resource_id, unit_code, status, version)
				VALUES ($1, 'flight_seat', $2, $3, 'available', 1)
			`, seatID, flightID, unitCode)

			if err != nil {
				return fmt.Errorf("insert seat %s: %w", unitCode, err)
			}
			seatNum++
		}
	}
	return nil
}
