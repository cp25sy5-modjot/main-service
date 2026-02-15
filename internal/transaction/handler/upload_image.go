package transactionhandler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	ext := "png"
	ct := c.Get("Content-Type")
	if strings.Contains(ct, "jpeg") || strings.Contains(ct, "jpg") {
		ext = "jpg"
	}

	ctx := c.Context()

	draftID := xid.New().String()

	path, err := h.storage.Save(ctx, userID, draftID, imageData, ext)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("Failed to store image: %v", err))
	}

	fullPath := filepath.Join("/uploads", path)

	if _, err := os.Stat(fullPath); err != nil {
		return fiber.NewError(500, "file not saved")
	}

	_, err = h.draftService.SaveDraft(ctx, draftID, userID, draft.NewDraftRequest{
		Path:      path,
		Items:     []draft.DraftItem{},
		CreatedAt: createAt,
	})
	if err != nil {
		return fiber.NewError(500, "failed to create draft")
	}
	
	task, err := tasks.NewBuildTransactionTask(userID, path, draftID)
	if err != nil {
		return fiber.NewError(500, "Failed to create job")
	}

	_, err = h.asynqClient.Enqueue(task,
		asynq.TaskID(draftID), // 🔥 กัน enqueue ซ้ำ
		asynq.MaxRetry(5),
		asynq.Timeout(10*time.Minute),
		asynq.ProcessIn(3*time.Second),
	)

	if err != nil {
		h.draftService.DeleteDraft(ctx, draftID, userID)
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
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to read uploaded image")
	}

	return data, nil
}
