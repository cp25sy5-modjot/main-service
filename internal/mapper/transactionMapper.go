package mapper

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func BuildTransactionResponse(tx *e.Transaction) *m.TransactionRes {
	items := BuildTransactionItemResponses(tx.Items)
	total := CalculateTotal(items)
	return &m.TransactionRes{
		TransactionID: tx.TransactionID,
		Title:         tx.Title,
		Date:          tx.Date,
		Type:          string(tx.Type),
		Total:         total,
		Items:         items,
	}
}

func BuildTransactionResponses(transactions []e.Transaction) []m.TransactionRes {
	if len(transactions) == 0 {
		return []m.TransactionRes{}
	}
	transactionResponses := make([]m.TransactionRes, 0, len(transactions))
	for i := range transactions {
		res := BuildTransactionResponse(&transactions[i])
		transactionResponses = append(transactionResponses, *res)
	}
	return transactionResponses
}

func BuildTransactionItemResponse(item *e.TransactionItem) *m.TransactionItemRes {
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

func BuildTransactionItemResponses(items []e.TransactionItem) []m.TransactionItemRes {
	if len(items) == 0 {
		return []m.TransactionItemRes{}
	}
	itemResponses := make([]m.TransactionItemRes, 0, len(items))
	for i := range items {
		res := BuildTransactionItemResponse(&items[i])
		itemResponses = append(itemResponses, *res)
	}
	return itemResponses
}

func ParseTransactionInsertReqToServiceInput(
	req *m.TransactionInsertReq,
) *m.TransactionCreateInput {
	return &m.TransactionCreateInput{
		Title: req.Title,
		Date:  req.Date,
		Items: MapTransactionItemReqToServiceInput(req.Items),
	}
}

func ParseTransactionUpdateReqToServiceInput(
	req *m.TransactionUpdateReq,
) *m.TransactionUpdateInput {
	return &m.TransactionUpdateInput{
		Title: req.Title,
		Date:  req.Date,
		Items: MapTransactionItemReqToServiceInput(req.Items),
	}
}

func MapTransactionItemReqToServiceInput(items []m.TransactionItemReq) []m.TransactionItemInput {
	if len(items) == 0 {
		return []m.TransactionItemInput{}
	}
	mappedItems := make([]m.TransactionItemInput, len(items))
	for i, item := range items {
		mappedItems[i] = m.TransactionItemInput{
			Title:      item.Title,
			Price:      item.Price,
			CategoryID: item.CategoryID,
		}
	}
	return mappedItems
}

func CalculateMonthTotal(transactions []m.TransactionRes) float64 {
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

func CalculateTotal(items []m.TransactionItemRes) float64 {
	if len(items) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range items {
		total += item.Price
	}
	return total
}

func ParseTransactionInsertReqToFavoriteItemCreateInput(
	userID string,
	req *m.TransactionInsertReq,
) *m.FavoriteItemCreateInput {
	return &m.FavoriteItemCreateInput{
		UserID:     userID,
		Title:      req.Items[0].Title,
		CategoryID: req.Items[0].CategoryID,
		Price:      req.Items[0].Price,
	}
}
