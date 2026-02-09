package transactionhandler

import (
	"strings"
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	draft "github.com/cp25sy5-modjot/main-service/internal/draft"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	fav "github.com/cp25sy5-modjot/main-service/internal/favorite_item/service"
)

type Handler struct {
	service      txsvc.Service // <- use interface, not *Service
	asynqClient  *asynq.Client
	storage      storage.Storage
	draftService draft.Service
	favService   fav.Service
}

func NewHandler(
	svc txsvc.Service, 
	client *asynq.Client, 
	st storage.Storage, 
	draftSvc draft.Service,
	favSvc fav.Service) *Handler {
	return &Handler{
		service:      svc,
		asynqClient:  client,
		storage:      st,
		draftService: draftSvc,
		favService:   favSvc,
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

	var input = mapper.ParseTransactionInsertReqToServiceInput(&req)

	resp, err := h.service.Create(userID, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	isNewFavorite := req.IsCreateNewFavorite
	if isNewFavorite {
		h.favService.Create(mapper.ParseTransactionInsertReqToFavoriteItemCreateInput(userID, &req))
	}
	return sresp.Created(c, mapper.BuildTransactionResponse(resp), "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	date := c.Query("date")

	categoryStr := c.Query("category")
	var categories []string
	if categoryStr != "" {
		categories = strings.Split(categoryStr, ",")
	}

	filter := &m.TransactionFilter{
		Date:       utils.ConvertStringToTime(date),
		Categories: categories,
	}

	months, err := h.service.GetAllComparePreviousMonthAndByUserIDWithFilter(userID, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve transactions")
	}

	currMonth := mapper.BuildTransactionResponses(months.CurrentMonth)
	previousMonth := mapper.BuildTransactionResponses(months.PreviousMonth)

	resp := m.TransactionCompareMonthResponse{
		Transactions:          currMonth,
		CurrentMonthTotal:     mapper.CalculateMonthTotal(currMonth),
		PreviousMonthTotal:    mapper.CalculateMonthTotal(previousMonth),
		CurrentMonthItemCount: months.CurrentMonthItemCount,
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

	return sresp.OK(c, mapper.BuildTransactionResponse(resp), "Transaction retrieved successfully")
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
		date := time.Now().UTC()
		req.Date = &date
	}

	input := mapper.ParseTransactionUpdateReqToServiceInput(&req)

	resp, err := h.service.Update(TransactionSearchParams, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update transaction")
	}
	return sresp.OK(c, mapper.BuildTransactionResponse(resp), "Transaction updated successfully")
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
