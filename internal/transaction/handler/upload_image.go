package transactionhandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/rs/xid"

	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
)

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
		return fiber.NewError(
			fiber.StatusInternalServerError,
			fmt.Sprintf("Failed to store image: %v", err),
		)
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

func getImageData(c *fiber.Ctx) ([]byte, error) {
	image, err := c.FormFile("image")
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Failed to upload image")
	}

	contentType := image.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Uploaded file is not a valid image")
	}

	file, err := image.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to process uploaded image")
	}
	defer file.Close()

	imageData := make([]byte, image.Size)
	_, err = file.Read(imageData)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to read uploaded image")
	}

	return imageData, nil
}
