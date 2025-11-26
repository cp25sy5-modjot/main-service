package transactionhandler

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/rs/xid"

	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
)

// getImageData is assumed to be your existing helper
// func getImageData(c *fiber.Ctx) ([]byte, error) { ... }

func (h *Handler) UploadImage(c *fiber.Ctx) error {
	imageData, err := getImageData(c)
	if err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	// Determine extension
	ext := "png"
	ct := c.Get("Content-Type")
	if strings.Contains(ct, "jpeg") || strings.Contains(ct, "jpg") {
		ext = "jpg"
	}

	ctx := context.Background()

	// 1. Save file to storage
	path, err := h.storage.Save(ctx, userID, imageData, ext)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to store image")
	}

	traceID := xid.New().String()

	// 2. Build async job
	task, err := tasks.NewBuildTransactionTask(userID, path, traceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create job")
	}

	// 3. Enqueue
	info, err := h.asynqClient.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(10*time.Minute),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to enqueue job")
	}

	return r.OK(c, fiber.Map{
		"job_id":   info.ID,
		"status":   "queued",
		"trace_id": traceID,
		"path":     path,
	}, "Image uploaded. Transaction will be processed asynchronously.")
}
