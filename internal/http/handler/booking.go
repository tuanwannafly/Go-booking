package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/service"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idempotencyKey := c.GetString("idempotency_key")
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(), userID.(uuid.UUID), req, idempotencyKey)
	if err != nil {
		if err.Error() == "idempotency key conflict: different request body" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Idempotency-Key conflict: different request body"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if this was a cached response
	if cached, ok := c.Get("idempotency_cached"); ok && cached == true {
		c.Header("X-Idempotency-Cached", "true")
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	booking, err := h.bookingService.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// Check ownership
	if booking.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to view this booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) GetUserBookings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		// parse page
	}
	if ps := c.Query("page_size"); ps != "" {
		// parse page_size
	}

	bookings, total, err := h.bookingService.GetUserBookings(c.Request.Context(), userID.(uuid.UUID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.PaginatedResponse{
		Data:       bookings,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	var req domain.ConfirmBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.bookingService.ConfirmBooking(c.Request.Context(), bookingID, req)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		if err.Error() == "booking is not in pending status" {
			c.JSON(http.StatusConflict, gin.H{"error": "booking cannot be confirmed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	var req domain.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.bookingService.CancelBooking(c.Request.Context(), bookingID, req)
	if err != nil {
		if err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		if err.Error() == "not authorized to cancel this booking" {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to cancel this booking"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// IdempotencyMiddleware handles idempotency key processing
func IdempotencyMiddleware(bookingService *service.BookingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to mutating methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.Next()
			return
		}

		// Read and store request body for hash comparison
		body, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		// Check if key exists
		existingBooking, err := bookingService.GetBookingByIdempotencyKey(c.Request.Context(), idempotencyKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		_ = existingBooking

		// Store body for later use
		c.Set("idempotency_key", idempotencyKey)
		c.Set("idempotency_body", body)
		c.Next()
	}
}

func computeRequestHash(body []byte) string {
	// Simple hash - in production use proper serialization
	return string(body) // placeholder
}
