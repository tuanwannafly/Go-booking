package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

type FlightRepository struct {
	db *pgxpool.Pool
}

func NewFlightRepository(db *pgxpool.Pool) *FlightRepository {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) Create(ctx context.Context, flight *domain.Flight) error {
	query := `
		INSERT INTO flights (id, code, origin, destination, departure_time, arrival_time, base_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query, flight.ID, flight.Code, flight.Origin, flight.Destination, flight.DepartureTime, flight.ArrivalTime, flight.BasePrice, flight.CreatedAt, flight.UpdatedAt)
	return err
}

func (r *FlightRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Flight, error) {
	query := `SELECT id, code, origin, destination, departure_time, arrival_time, base_price, created_at, updated_at FROM flights WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var flight domain.Flight
	err := row.Scan(&flight.ID, &flight.Code, &flight.Origin, &flight.Destination, &flight.DepartureTime, &flight.ArrivalTime, &flight.BasePrice, &flight.CreatedAt, &flight.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &flight, err
}

func (r *FlightRepository) Search(ctx context.Context, params domain.FlightSearchParams) ([]domain.Flight, int64, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", params.Date)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid date format: %w", err)
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	countQuery := `
		SELECT COUNT(*) FROM flights
		WHERE origin = $1 AND destination = $2 AND departure_time >= $3 AND departure_time < $4
	`
	var total int64
	err = r.db.QueryRow(ctx, countQuery, params.Origin, params.Destination, startOfDay, endOfDay).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Flight{}, 0, nil
	}

	offset := (params.Page - 1) * params.PageSize
	query := `
		SELECT id, code, origin, destination, departure_time, arrival_time, base_price, created_at, updated_at
		FROM flights
		WHERE origin = $1 AND destination = $2 AND departure_time >= $3 AND departure_time < $4
		ORDER BY departure_time
		LIMIT $5 OFFSET $6
	`
	rows, err := r.db.Query(ctx, query, params.Origin, params.Destination, startOfDay, endOfDay, params.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var flights []domain.Flight
	for rows.Next() {
		var flight domain.Flight
		err := rows.Scan(&flight.ID, &flight.Code, &flight.Origin, &flight.Destination, &flight.DepartureTime, &flight.ArrivalTime, &flight.BasePrice, &flight.CreatedAt, &flight.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		flights = append(flights, flight)
	}

	return flights, total, nil
}

type HotelRepository struct {
	db *pgxpool.Pool
}

func NewHotelRepository(db *pgxpool.Pool) *HotelRepository {
	return &HotelRepository{db: db}
}

func (r *HotelRepository) Create(ctx context.Context, hotel *domain.Hotel) error {
	query := `
		INSERT INTO hotels (id, name, city, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, hotel.ID, hotel.Name, hotel.City, hotel.Address, hotel.CreatedAt, hotel.UpdatedAt)
	return err
}

func (r *HotelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Hotel, error) {
	query := `SELECT id, name, city, address, created_at, updated_at FROM hotels WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var hotel domain.Hotel
	err := row.Scan(&hotel.ID, &hotel.Name, &hotel.City, &hotel.Address, &hotel.CreatedAt, &hotel.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &hotel, err
}

func (r *HotelRepository) Search(ctx context.Context, params domain.HotelSearchParams) ([]domain.Hotel, int64, error) {
	_, err := time.Parse("2006-01-02", params.CheckIn)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid checkin date: %w", err)
	}
	_, err = time.Parse("2006-01-02", params.CheckOut)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid checkout date: %w", err)
	}

	countQuery := `SELECT COUNT(*) FROM hotels WHERE city = $1`
	var total int64
	err = r.db.QueryRow(ctx, countQuery, params.City).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []domain.Hotel{}, 0, nil
	}

	offset := (params.Page - 1) * params.PageSize
	query := `
		SELECT id, name, city, address, created_at, updated_at
		FROM hotels
		WHERE city = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, params.City, params.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var hotels []domain.Hotel
	for rows.Next() {
		var hotel domain.Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.City, &hotel.Address, &hotel.CreatedAt, &hotel.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		hotels = append(hotels, hotel)
	}

	return hotels, total, nil
}

type RoomTypeRepository struct {
	db *pgxpool.Pool
}

func NewRoomTypeRepository(db *pgxpool.Pool) *RoomTypeRepository {
	return &RoomTypeRepository{db: db}
}

func (r *RoomTypeRepository) Create(ctx context.Context, rt *domain.RoomType) error {
	query := `
		INSERT INTO room_types (id, hotel_id, name, base_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, rt.ID, rt.HotelID, rt.Name, rt.BasePrice, rt.CreatedAt, rt.UpdatedAt)
	return err
}

func (r *RoomTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RoomType, error) {
	query := `SELECT id, hotel_id, name, base_price, created_at, updated_at FROM room_types WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var rt domain.RoomType
	err := row.Scan(&rt.ID, &rt.HotelID, &rt.Name, &rt.BasePrice, &rt.CreatedAt, &rt.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rt, err
}

func (r *RoomTypeRepository) GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]domain.RoomType, error) {
	query := `SELECT id, hotel_id, name, base_price, created_at, updated_at FROM room_types WHERE hotel_id = $1`
	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomTypes []domain.RoomType
	for rows.Next() {
		var rt domain.RoomType
		err := rows.Scan(&rt.ID, &rt.HotelID, &rt.Name, &rt.BasePrice, &rt.CreatedAt, &rt.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roomTypes = append(roomTypes, rt)
	}
	return roomTypes, nil
}
