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
	currMonth := buildTransactionResponses(months.CurrentMonth)
	previousMonth := buildTransactionResponses(months.PreviousMonth)
	resp := m.TransactionCompareMonthResponse{
		Transactions:       currMonth,
		CurrentMonthTotal:  calculateTotal(currMonth),
		PreviousMonthTotal: calculateTotal(previousMonth),
	}

	return sresp.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id
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

// PUT /transactions/:transaction_id
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

// DELETE /transactions/:transaction_id
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
func getTxIDAndProdID(c *fiber.Ctx) (string, error) {
	tx_id := c.Params("transaction_id")
	if tx_id == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "transaction_id and item_id parameters are required")
	}
	return tx_id, nil
}

func createTransactionSearchParams(c *fiber.Ctx) (*m.TransactionSearchParams, error) {
	tx_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &m.TransactionSearchParams{
		TransactionID: tx_id,
		UserID:        userID,
	}, nil
}

func buildTransactionResponse(tx *e.Transaction) *m.TransactionRes {
	return &m.TransactionRes{
		TransactionID: tx.TransactionID,
		Date:          tx.Date,
		Type:          string(tx.Type),
		Items:         buildTransactionItemResponses(tx.Items),
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

func buildTransactionItemResponse(item *e.TransactionItem) *m.TransactionItemRes {
	return &m.TransactionItemRes{
		TransactionID:     item.TransactionID,
		ItemID:            item.ItemID,
		Title:             item.Title,
		Price:             item.Price,
		CategoryID:        item.CategoryID,
		CategoryName:      item.Category.CategoryName,
		CategoryColorCode: item.Category.ColorCode,
	}
}

func buildTransactionItemResponses(items []e.TransactionItem) []m.TransactionItemRes {
	if len(items) == 0 {
		return []m.TransactionItemRes{}
	}
	itemResponses := make([]m.TransactionItemRes, 0, len(items))
	for _, item := range items {
		res := buildTransactionItemResponse(&item)
		itemResponses = append(itemResponses, *res)
	}
	return itemResponses
}

func parseTransactionInsertReqToServiceInput(
	req *m.TransactionInsertReq,
) *txsvc.TransactionCreateInput {
	return &txsvc.TransactionCreateInput{
		Date:  req.Date,
		Items: mapTransactionItemReqToServiceInput(req.Items),
	}
}

func parseTransactionUpdateReqToServiceInput(
	req *m.TransactionUpdateReq,
) *txsvc.TransactionUpdateInput {
	return &txsvc.TransactionUpdateInput{
		Date:  req.Date,
		Items: mapTransactionItemReqToServiceInput(req.Items),
	}
}

func mapTransactionItemReqToServiceInput(items []m.TransactionItemReq) []txsvc.TransactionItemInput {
	if len(items) == 0 {
		return []txsvc.TransactionItemInput{}
	}
	mappedItems := make([]txsvc.TransactionItemInput, len(items))
	for i, item := range items {
		mappedItems[i] = txsvc.TransactionItemInput{
			Title:      item.Title,
			Price:      item.Price,
			CategoryID: item.CategoryID,
		}
	}
	return mappedItems
}

func calculateTotal(transactions []m.TransactionRes) float64 {
	if len(transactions) == 0 {
		return 0
	}
	total := 0.0
	for _, tx := range transactions {
		for _, item := range tx.Items {
			total += item.Price
		}
	}
	return total
}
