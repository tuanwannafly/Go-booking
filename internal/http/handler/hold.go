package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gobooking/internal/domain"
	"gobooking/internal/repository/postgres"
	"gobooking/internal/service"
)

type HoldHandler struct {
	holdService *service.HoldService
}

func NewHoldHandler(holdService *service.HoldService) *HoldHandler {
	return &HoldHandler{holdService: holdService}
}

func (h *HoldHandler) HoldSeat(c *gin.Context) {
	seatID, err := uuid.Parse(c.Param("seatId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seat ID"})
		return
	}

	var req domain.HoldSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override with path parameter
	req.SeatID = seatID

	unit, err := h.holdService.HoldSeat(c.Request.Context(), req.SeatID, req.HoldDuration)
	if err != nil {
		if errors.Is(err, postgres.ErrInventoryNotAvailable) || errors.Is(err, postgres.ErrInventoryNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": "seat not available"})
			return
		}
		if errors.Is(err, service.ErrInvalidHoldDuration) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrHoldLockUnavailable) {
			c.JSON(http.StatusConflict, gin.H{"error": "seat is being held by another request"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory_unit": unit,
		"held_until":     unit.HeldUntil,
		"expires_in":     expiresInSeconds(unit.HeldUntil),
	})
}

func (h *HoldHandler) HoldRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	var req domain.HoldRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.RoomID = roomID

	unit, err := h.holdService.HoldRoom(c.Request.Context(), req.RoomID, req.HoldDuration)
	if err != nil {
		if errors.Is(err, postgres.ErrInventoryNotAvailable) || errors.Is(err, postgres.ErrInventoryNotFound) {
			c.JSON(http.StatusConflict, gin.H{"error": "room not available"})
			return
		}
		if errors.Is(err, service.ErrInvalidHoldDuration) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, postgres.ErrOptimisticLockConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "room temporarily unavailable, please retry"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory_unit": unit,
		"held_until":     unit.HeldUntil,
		"expires_in":     expiresInSeconds(unit.HeldUntil),
	})
}

func expiresInSeconds(deadline *time.Time) int {
	if deadline == nil {
		return 0
	}
	remaining := time.Until(*deadline)
	if remaining <= 0 {
		return 0
	}
	return int(remaining.Round(time.Second).Seconds())
}

func (h *HoldHandler) ReleaseHold(c *gin.Context) {
	inventoryUnitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inventory unit ID"})
		return
	}

	err = h.holdService.ReleaseHold(c.Request.Context(), inventoryUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "hold released"})
}
