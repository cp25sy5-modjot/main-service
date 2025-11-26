package transactionhandler

import (
	"strings"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *txsvc.Service
}

func NewHandler(service *txsvc.Service) *Handler {
	return &Handler{service}
}

// POST /transactions/manual
func (h *Handler) Create(c *fiber.Ctx) error {
	var req m.TransactionInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var input = parseTransactionInsertReqToServiceInput(&req)

	resp, err := h.service.Create(userID, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return sresp.Created(c, buildTransactionResponse(resp), "Transaction created successfully")
}

// POST /transactions/upload
func (h *Handler) UploadImage(c *fiber.Ctx) error {
	imageData, err := getImageData(c)
	if err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	resp, err := h.service.ProcessUploadedFile(imageData, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process the uploaded file")
	}

	return sresp.Created(c, buildTransactionResponse(resp), "File uploaded and processed successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	date := c.Query("date")
	filter := &m.TransactionFilter{
		Date: utils.ConvertStringToTime(date),
	}

	resp, err := h.service.GetAllByUserIDWithFilter(userID, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve transactions")
	}
	return sresp.OK(c, buildTransactionResponses(resp), "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id/product/:item_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetByID(TransactionSearchParams)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}

	return sresp.OK(c, buildTransactionResponse(resp), "Transaction retrieved successfully")
}

// PUT /transactions/:transaction_id/product/:item_id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req m.TransactionUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}

	input := parseTransactionUpdateReqToServiceInput(&req)

	resp, err := h.service.Update(TransactionSearchParams, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update transaction")
	}
	return sresp.OK(c, buildTransactionResponse(resp), "Transaction updated successfully")
}

// DELETE /transactions/:transaction_id/product/:item_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	TransactionSearchParams, err := createTransactionSearchParams(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(TransactionSearchParams); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete transaction")
	}

	return sresp.OK(c, nil, "Transaction deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	item_id := c.Params("item_id")
	if tx_id == "" || item_id == "" {
		return "", "", fiber.NewError(fiber.StatusBadRequest, "transaction_id and item_id parameters are required")
	}
	return tx_id, item_id, nil
}

func createTransactionSearchParams(c *fiber.Ctx) (*m.TransactionSearchParams, error) {
	tx_id, item_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &m.TransactionSearchParams{
		TransactionID: tx_id,
		ItemID:        item_id,
		UserID:        userID,
	}, nil
}

func buildTransactionResponse(tx *e.Transaction) *m.TransactionRes {
	return &m.TransactionRes{
		TransactionID:     tx.TransactionID,
		ItemID:            tx.ItemID,
		Title:             tx.Title,
		Price:             tx.Price,
		Quantity:          tx.Quantity,
		TotalPrice:        tx.Price * tx.Quantity,
		Date:              utils.ToUserLocal(tx.Date, ""),
		Type:              tx.Type,
		CategoryID:        tx.CategoryID,
		CategoryName:      tx.Category.CategoryName,
		CategoryColorCode: tx.Category.ColorCode,
	}
}

func buildTransactionResponses(transactions []e.Transaction) []m.TransactionRes {
	transactionResponses := make([]m.TransactionRes, 0, len(transactions))
	for _, tx := range transactions {
		res := buildTransactionResponse(&tx)
		transactionResponses = append(transactionResponses, *res)
	}
	return transactionResponses
}

func parseTransactionInsertReqToServiceInput(req *m.TransactionInsertReq) *txsvc.TransactionCreateInput {
	if req.Date.IsZero() {
		req.Date = utils.NowUTC()
	}
	return &txsvc.TransactionCreateInput{
		Title:      req.Title,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Date:       utils.NormalizeToUTC(req.Date, ""),
		CategoryID: req.CategoryID,
	}
}

func parseTransactionUpdateReqToServiceInput(req *m.TransactionUpdateReq) *txsvc.TransactionUpdateInput {
	if req.Date.IsZero() {
		req.Date = utils.NowUTC()
	}
	return &txsvc.TransactionUpdateInput{
		Title:      req.Title,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Date:       utils.NormalizeToUTC(req.Date, ""),
		CategoryID: req.CategoryID,
	}
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
