package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/repository/postgres"
)

type SearchService struct {
	flightRepo    *postgres.FlightRepository
	hotelRepo     *postgres.HotelRepository
	roomTypeRepo  *postgres.RoomTypeRepository
	inventoryRepo *postgres.InventoryRepository
}

func NewSearchService(
	flightRepo *postgres.FlightRepository,
	hotelRepo *postgres.HotelRepository,
	roomTypeRepo *postgres.RoomTypeRepository,
	inventoryRepo *postgres.InventoryRepository,
) *SearchService {
	return &SearchService{
		flightRepo:    flightRepo,
		hotelRepo:     hotelRepo,
		roomTypeRepo:  roomTypeRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *SearchService) SearchFlights(ctx context.Context, params domain.FlightSearchParams) ([]domain.Flight, int64, error) {
	return s.flightRepo.Search(ctx, params)
}

func (s *SearchService) SearchHotels(ctx context.Context, params domain.HotelSearchParams) ([]domain.Hotel, int64, error) {
	hotels, total, err := s.hotelRepo.Search(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// For each hotel, get room types with available inventory
	for i := range hotels {
		roomTypes, err := s.roomTypeRepo.GetByHotelID(ctx, hotels[i].ID)
		if err != nil {
			return nil, 0, err
		}

		// Check availability for each room type
		for j := range roomTypes {
			available, err := s.inventoryRepo.GetAvailableByResource(ctx, domain.ResourceTypeHotelRoom, roomTypes[j].ID)
			if err != nil {
				return nil, 0, err
			}
			roomTypes[j].AvailableRooms = len(available)
		}
		hotels[i].RoomTypes = roomTypes
	}

	return hotels, total, nil
}

func (s *SearchService) GetFlightWithAvailability(ctx context.Context, flightID uuid.UUID) (*domain.Flight, error) {
	flight, err := s.flightRepo.GetByID(ctx, flightID)
	if err != nil {
		return nil, err
	}

	if flight == nil {
		return nil, fmt.Errorf("flight not found")
	}

	// Get available seats
	seats, err := s.inventoryRepo.GetAvailableByResource(ctx, domain.ResourceTypeFlightSeat, flightID)
	if err != nil {
		return nil, err
	}
	flight.AvailableSeats = len(seats)

	return flight, nil
}

func (s *SearchService) GetHotelWithAvailability(ctx context.Context, hotelID uuid.UUID, checkIn, checkOut time.Time) (*domain.Hotel, error) {
	hotel, err := s.hotelRepo.GetByID(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	if hotel == nil {
		return nil, fmt.Errorf("hotel not found")
	}

	roomTypes, err := s.roomTypeRepo.GetByHotelID(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	for i := range roomTypes {
		// For hotels, we need to check availability across the date range
		// This is simplified - in reality you'd check each night
		available, err := s.inventoryRepo.GetAvailableByResource(ctx, domain.ResourceTypeHotelRoom, roomTypes[i].ID)
		if err != nil {
			return nil, err
		}
		roomTypes[i].AvailableRooms = len(available)
	}

	hotel.RoomTypes = roomTypes
	return hotel, nil
}
