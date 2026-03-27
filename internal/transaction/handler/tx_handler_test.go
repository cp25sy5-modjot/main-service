package transactionhandler_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	draftmocks "github.com/cp25sy5-modjot/main-service/internal/draft/mocks"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	queuemocks "github.com/cp25sy5-modjot/main-service/internal/queue/mocks"
	storagemocks "github.com/cp25sy5-modjot/main-service/internal/storage/mocks"
	transactionhandler "github.com/cp25sy5-modjot/main-service/internal/transaction/handler"
)

func TestTransactionHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	t.Run("Upload Image Success", func(t *testing.T) {
		app := fiber.New()

		mockStorage := storagemocks.NewMockStorage(t)
		mockDraft := draftmocks.NewMockService(t)
		mockQueue := queuemocks.NewMockQueue(t)

		h := transactionhandler.NewHandler(
			nil,
			mockQueue,
			mockStorage,
			mockDraft,
			nil,
		)

		app.Post("/upload", h.UploadImage)

		// storage
		mockStorage.EXPECT().
			Save(
				mock.Anything,
				"u-1",
				mock.Anything,
				mock.Anything,
				"png",
			).
			Return("test.png", nil)

		// draft
		mockDraft.EXPECT().
			SaveDraft(
				mock.Anything,
				mock.Anything,
				"u-1",
				mock.Anything,
			).
			Return(nil, nil)

		// queue
		mockQueue.EXPECT().
			Enqueue(
				mock.Anything, // task
				mock.Anything, // TaskID
				mock.Anything, // MaxRetry
				mock.Anything, // Timeout
				mock.Anything, // ProcessIn
			).
			Return(&asynq.TaskInfo{}, nil)

		//FIX: multipart with Content-Type
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		partHeader := make(textproto.MIMEHeader)
		partHeader.Set("Content-Disposition", `form-data; name="image"; filename="test.png"`)
		partHeader.Set("Content-Type", "image/png") 

		part, _ := writer.CreatePart(partHeader)
		part.Write([]byte("fake-image"))

		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
