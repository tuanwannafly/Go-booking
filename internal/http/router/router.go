package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gobooking/internal/config"
	"gobooking/internal/domain"
	"gobooking/internal/http/handler"
	"gobooking/internal/http/middleware"
	"gobooking/internal/repository/postgres"
	"gobooking/internal/repository/redis"
	"gobooking/internal/service"
)

type Router struct {
	engine *gin.Engine
	cfg    *config.Config
	db     *pgxpool.Pool
	redis  *redis.RedisClient
}

func NewRouter(cfg *config.Config, db *pgxpool.Pool, redisClient *redis.RedisClient) *Router {
	r := &Router{
		engine: gin.Default(),
		cfg:    cfg,
		db:     db,
		redis:  redisClient,
	}
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// Middleware
	r.engine.Use(middleware.RequestIDMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())
	r.engine.Use(middleware.RecoveryMiddleware())
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.RateLimitMiddleware(100, 200))

	// Health checks
	r.engine.GET("/healthz", r.healthCheck)
	r.engine.GET("/readyz", r.readyCheck)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(r.db)
	flightRepo := postgres.NewFlightRepository(r.db)
	hotelRepo := postgres.NewHotelRepository(r.db)
	roomTypeRepo := postgres.NewRoomTypeRepository(r.db)
	inventoryRepo := postgres.NewInventoryRepository(r.db)
	bookingRepo := postgres.NewBookingRepository(r.db)
	bookingItemRepo := postgres.NewBookingItemRepository(r.db)
	idempotencyRepo := postgres.NewIdempotencyRepository(r.db)
	paymentRepo := postgres.NewPaymentRepository(r.db)

	// Initialize services
	searchService := service.NewSearchService(flightRepo, hotelRepo, roomTypeRepo, inventoryRepo)
	holdService := service.NewHoldService(inventoryRepo, r.redis, r.cfg.Worker.HoldDurationMinutes)
	bookingService := service.NewBookingService(
		bookingRepo,
		bookingItemRepo,
		inventoryRepo,
		idempotencyRepo,
		paymentRepo,
		userRepo,
		holdService,
		domain.DefaultRefundPolicy,
		r.cfg.Worker.BookingExpireMinutes,
	)

	// Initialize handlers
	searchHandler := handler.NewSearchHandler(searchService)
	holdHandler := handler.NewHoldHandler(holdService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// API routes
	api := r.engine.Group("/api/v1")
	{
		// Search routes (public)
		api.GET("/flights/search", searchHandler.SearchFlights)
		api.GET("/flights/:id", searchHandler.GetFlight)
		api.GET("/hotels/search", searchHandler.SearchHotels)
		api.GET("/hotels/:id", searchHandler.GetHotel)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(r.cfg.App.JWTSecret))
		{
			// Hold routes
			protected.POST("/flights/:id/seats/:seatId/hold", holdHandler.HoldSeat)
			protected.POST("/hotels/rooms/:roomId/hold", holdHandler.HoldRoom)
			protected.DELETE("/holds/:id", holdHandler.ReleaseHold)

			// Booking routes
			protected.POST("/bookings", bookingHandler.CreateBooking)
			protected.GET("/bookings", bookingHandler.GetUserBookings)
			protected.GET("/bookings/:id", bookingHandler.GetBooking)
			protected.POST("/bookings/:id/confirm", bookingHandler.ConfirmBooking)
			protected.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)
		}
	}

	// Swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   "2024-01-01T00:00:00Z",
	})
}

func (r *Router) readyCheck(c *gin.Context) {
	// Check database
	ctx := c.Request.Context()
	if err := r.db.Ping(ctx); err != nil {
		c.JSON(503, gin.H{"status": "not ready", "database": "unhealthy"})
		return
	}

	// Check Redis
	if err := r.redis.Ping(ctx); err != nil {
		c.JSON(503, gin.H{"status": "not ready", "redis": "unhealthy"})
		return
	}

	c.JSON(200, gin.H{"status": "ready"})
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
