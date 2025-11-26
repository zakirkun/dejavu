package handler

import (
	"github.com/dejavu/backend/internal/domain"
	"github.com/dejavu/backend/internal/repository"
	"github.com/dejavu/backend/internal/service"
	"github.com/dejavu/backend/pkg/database"
	"github.com/dejavu/backend/pkg/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type DeployHandler struct {
	service *service.DeploymentService
	queue   *queue.Queue
}

func NewDeployHandler(db *database.DB, nats *queue.Queue) *DeployHandler {
	deployRepo := repository.NewDeploymentRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	deployService := service.NewDeploymentService(deployRepo, projectRepo, nats)
	return &DeployHandler{
		service: deployService,
		queue:   nats,
	}
}

func (h *DeployHandler) Trigger(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req domain.TriggerDeployRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	deployment, err := h.service.Trigger(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(deployment)
}

func (h *DeployHandler) GetStatus(c *fiber.Ctx) error {
	deployID := c.Params("id")

	deployment, err := h.service.GetStatus(deployID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(deployment)
}

func (h *DeployHandler) StreamLogs(c *fiber.Ctx) error {
	// Upgrade to WebSocket
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return websocket.New(func(ws *websocket.Conn) {
			deployID := c.Params("id")

			// Subscribe to logs for this deployment
			sub, err := h.queue.Subscribe("BUILDS.logs."+deployID, func(data []byte) {
				ws.WriteMessage(websocket.TextMessage, data)
			})
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("Error subscribing to logs"))
				return
			}
			defer sub.Unsubscribe()

			// Keep connection alive
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					break
				}
			}
		})(c)
	}

	return c.Status(fiber.StatusUpgradeRequired).JSON(fiber.Map{
		"error": "WebSocket upgrade required",
	})
}

