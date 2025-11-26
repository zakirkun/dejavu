package handler

import (
	"github.com/dejavu/backend/internal/repository"
	"github.com/dejavu/backend/internal/service"
	"github.com/dejavu/backend/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type BillingHandler struct {
	service *service.BillingService
}

func NewBillingHandler(db *database.DB) *BillingHandler {
	repo := repository.NewBillingRepository(db)
	billingService := service.NewBillingService(repo)
	return &BillingHandler{service: billingService}
}

func (h *BillingHandler) GetBalance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"balance": balance,
	})
}

func (h *BillingHandler) GetUsageHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	history, err := h.service.GetUsageHistory(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(history)
}

func (h *BillingHandler) AddCredits(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.service.AddCredits(userID, req.Amount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Credits added successfully",
	})
}

