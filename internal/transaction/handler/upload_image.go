package transactionhandler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/rs/xid"

	draft "github.com/cp25sy5-modjot/main-service/internal/draft"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
)

func (h *Handler) UploadImage(c *fiber.Ctx) error {

	createAt := time.Now()
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

	draftID := xid.New().String()

	_, err = h.draftService.SaveDraft(ctx, draftID, userID, draft.NewDraftRequest{
		// Title: "Slip Image Upload",
		// Date:  time.Now(),
		Path:      path,
		Items:     []draft.DraftItem{},
		CreatedAt: createAt,
	})

	if err != nil {
		return fiber.NewError(500, "failed to create draft")
	}

	// 3. Build async job
	task, err := tasks.NewBuildTransactionTask(userID, path, draftID)
	if err != nil {
		return fiber.NewError(500, "Failed to create job")
	}

	// 4. Enqueue
	_, err = h.asynqClient.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Timeout(10*time.Minute),
	)

	if err != nil {
		// rollback draft ถ้า enqueue ไม่ผ่าน
		h.draftService.DeleteDraft(ctx, draftID)

		return fiber.NewError(500, "Failed to enqueue job")
	}

	return r.OK(c, fiber.Map{
		"draft_id": draftID,
		"status":   "queued",
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
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	imageData := make([]byte, image.Size)
	_, err = file.Read(imageData)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to read uploaded image")
	}

	return imageData, nil
}
