package handler

import (
	"binancetrading/internal/interfaces"
	"github.com/gofiber/fiber/v2"
	"time"
)

// APIHandler manages API endpoints
type APIHandler struct {
	storage interfaces.Storage
}

func NewAPIHandler(storage interfaces.Storage) *APIHandler {
	return &APIHandler{storage: storage}
}

// SetupRoutes sets up the API routes
func (h *APIHandler) SetupRoutes(app *fiber.App) {
	app.Get("/candlesticks/:symbol", h.GetCandlesticks)
}

// GetCandlesticks handles GET requests to retrieve candlesticks
func (h *APIHandler) GetCandlesticks(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if startTime == "" {
		startTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}
	if endTime == "" {
		endTime = time.Now().Format(time.RFC3339)
	}

	candlesticks, err := h.storage.GetCandlesticks(symbol, startTime, endTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(candlesticks)
}
