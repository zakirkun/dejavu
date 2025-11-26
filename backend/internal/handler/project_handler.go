package handler

import (
	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/internal/repository"
	"github.com/dejavu/backend/internal/service"
	"github.com/dejavu/backend/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(db *database.DB) *ProjectHandler {
	repo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(repo)
	return &ProjectHandler{service: projectService}
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req domain.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	project, err := h.service.Create(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func (h *ProjectHandler) Get(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	projectID := c.Params("id")

	project, err := h.service.GetByID(projectID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(project)
}

func (h *ProjectHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	projects, err := h.service.List(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(projects)
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	projectID := c.Params("id")

	var req domain.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.service.Update(projectID, userID, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Project updated successfully",
	})
}

func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	projectID := c.Params("id")

	if err := h.service.Delete(projectID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}

