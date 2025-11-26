package handler

import (
	"os"

	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/internal/repository"
	"github.com/dejavu/backend/internal/service"
	"github.com/dejavu/backend/pkg/cache"
	"github.com/dejavu/backend/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(db *database.DB, redis *cache.Cache) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	authService := service.NewAuthService(userRepo, jwtSecret)
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	user, err := h.service.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	response, err := h.service.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

