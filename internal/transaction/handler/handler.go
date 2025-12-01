package transactionhandler

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

type Handler struct {
	service     txsvc.Service // <- use interface, not *Service
	asynqClient *asynq.Client
	storage     storage.Storage
}

func NewHandler(svc txsvc.Service, client *asynq.Client, st storage.Storage) *Handler {
	return &Handler{
		service:     svc,
		asynqClient: client,
		storage:     st,
	}
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

	months, err := h.service.GetAllComparePreviousMonthAndByUserIDWithFilter(userID, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve transactions")
	}
	resp := m.TransactionCompareMonthResponse{
		Transactions:       months.CurrentMonth,
		CurrentMonthTotal:  calculateTotal(months.CurrentMonth),
		PreviousMonthTotal: calculateTotal(months.PreviousMonth),
	}

	return sresp.OK(c, resp, "Transactions retrieved successfully")
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

	return sresp.OK(c, resp, "Transaction retrieved successfully")
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

	if req.Date == nil {
		date := time.Now()
		req.Date = &date
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
		TransactionID: tx.TransactionID,
		ItemID:        tx.ItemID,
		Title:         tx.Title,
		Price:         tx.Price,
		Date:          tx.Date,
		Type:          tx.Type,
		CategoryID:    tx.CategoryID,
	}
}

func buildTransactionResponses(transactions []e.Transaction) []m.TransactionRes {
	if len(transactions) == 0 {
		return []m.TransactionRes{}
	}
	transactionResponses := make([]m.TransactionRes, 0, len(transactions))
	for _, tx := range transactions {
		res := buildTransactionResponse(&tx)
		transactionResponses = append(transactionResponses, *res)
	}
	return transactionResponses
}

func parseTransactionInsertReqToServiceInput(req *m.TransactionInsertReq) *txsvc.TransactionCreateInput {
	return &txsvc.TransactionCreateInput{
		Title:      req.Title,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Date:       req.Date,
		CategoryID: req.CategoryID,
	}
}

func parseTransactionUpdateReqToServiceInput(req *m.TransactionUpdateReq) *txsvc.TransactionUpdateInput {
	return &txsvc.TransactionUpdateInput{
		Title:      req.Title,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Date:       req.Date,
		CategoryID: req.CategoryID,
	}
}

func calculateTotal(transactions []m.TransactionRes) float64 {
	if len(transactions) == 0 {
		return 0
	}
	total := 0.0
	for _, tx := range transactions {
		total += tx.Price
	}
	return total
}
