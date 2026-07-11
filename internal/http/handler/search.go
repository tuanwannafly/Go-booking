package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/service"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) SearchFlights(c *gin.Context) {
	var params domain.FlightSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	flights, total, err := h.searchService.SearchFlights(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.PaginatedResponse{
		Data:       flights,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: int((total + int64(params.PageSize) - 1) / int64(params.PageSize)),
	})
}

func (h *SearchHandler) SearchHotels(c *gin.Context) {
	var params domain.HotelSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// Validate dates
	checkIn, err := time.Parse("2006-01-02", params.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkin date format (YYYY-MM-DD)"})
		return
	}
	checkOut, err := time.Parse("2006-01-02", params.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkout date format (YYYY-MM-DD)"})
		return
	}
	if checkOut.Before(checkIn) || checkOut.Equal(checkIn) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "checkout must be after checkin"})
		return
	}

	hotels, total, err := h.searchService.SearchHotels(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.PaginatedResponse{
		Data:       hotels,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Total:      total,
		TotalPages: int((total + int64(params.PageSize) - 1) / int64(params.PageSize)),
	})
}

func (h *SearchHandler) GetFlight(c *gin.Context) {
	flightID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid flight ID"})
		return
	}

	flight, err := h.searchService.GetFlightWithAvailability(c.Request.Context(), flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if flight == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "flight not found"})
		return
	}

	c.JSON(http.StatusOK, flight)
}

func (h *SearchHandler) GetHotel(c *gin.Context) {
	hotelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hotel ID"})
		return
	}

	checkInStr := c.Query("checkin")
	checkOutStr := c.Query("checkout")

	var checkIn, checkOut time.Time
	if checkInStr != "" {
		checkIn, _ = time.Parse("2006-01-02", checkInStr)
	}
	if checkOutStr != "" {
		checkOut, _ = time.Parse("2006-01-02", checkOutStr)
	}

	hotel, err := h.searchService.GetHotelWithAvailability(c.Request.Context(), hotelID, checkIn, checkOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hotel == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hotel not found"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}
