package transaction

import (
	"log"
	"strings"

	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	successResp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	model "github.com/cp25sy5-modjot/main-service/internal/transaction/model"
	svc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *svc.Service
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{service}
}

// POST /transactions
func (h *Handler) Create(c *fiber.Ctx) error {
	var req model.TransactionInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	var tx model.Transaction
	_ = utils.MapStructs(req, &tx)
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	tx.UserID = userID
	tx.Type = "manual"
	createdTx, err := h.service.Create(&tx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var res model.TransactionRes
	utils.MapStructs(createdTx, &res)
	return successResp.Created(c, res, "Transaction created successfully")
}

func (h *Handler) UploadImage(c *fiber.Ctx) error {
	// Parse the uploaded image
	image, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to upload image")
	}

	contentType := image.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is not an image",
		})
	}

	file, err := image.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process uploaded image")
	}
	defer file.Close()

	imageData := make([]byte, image.Size)
	_, err = file.Read(imageData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read uploaded image")
	}

	// Get user ID from JWT claims
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	createdTx, err := h.service.ProcessUploadedFile(imageData, userID)
	if err != nil {
		log.Printf("Failed to process uploaded file: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process the uploaded file")
	}

	var res model.TransactionRes
	utils.MapStructs(createdTx, &res)
	return successResp.Created(c, res, "File uploaded and processed successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	transactions, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve transactions")
	}
	var resp []model.TransactionRes
	_ = utils.MapStructs(transactions, &resp)
	return successResp.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id/product/:product_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}
	transaction, err := h.service.GetByID(TransactionSearchParams)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}
	var resp model.TransactionRes
	_ = utils.MapStructs(transaction, &resp)
	return successResp.OK(c, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:transaction_id/product/:product_id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req model.TransactionUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}
	if err := h.service.Update(TransactionSearchParams, &req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update transaction")
	}
	var resp model.TransactionRes
	_ = utils.MapStructs(req, &resp)
	return successResp.OK(c, resp, "Transaction updated successfully")
}

// DELETE /transactions/:transaction_id/product/:product_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}
	if err := h.service.Delete(TransactionSearchParams); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete transaction")
	}
	return successResp.OK(c, nil, "Transaction deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	item_id := c.Params("product_id")
	if tx_id == "" || item_id == "" {
		return "", "", fiber.NewError(fiber.StatusBadRequest, "transaction_id and product_id parameters are required")
	}
	return tx_id, item_id, nil
}

func createTransactionSearchParams(c *fiber.Ctx) (*model.TransactionSearchParams, error) {
	tx_id, item_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &model.TransactionSearchParams{
		TransactionID: tx_id,
		ItemID:        item_id,
		UserID:        userID,
	}, nil
}
