package mapper

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func MapConfirmDraftToCreateInput(d *m.ConfirmRequest) *m.TransactionCreateInput {

	var items []m.TransactionItemInput

	for _, it := range d.Items {
		items = append(items, m.TransactionItemInput{
			Title:      it.Title,
			Price:      it.Price,
			CategoryID: it.CategoryID,
		})
	}

	return &m.TransactionCreateInput{
		Title: d.Title,
		Date:  *d.Date,
		Items: items,
	}
}

func BuildDraftRes(
	d m.DraftTxn,
	categoryMap map[string]e.Category,
) m.DraftRes {

	items := make([]m.DraftItemRes, 0, len(d.Items))

	for _, it := range d.Items {
		cat, ok := categoryMap[it.CategoryID]

		item := m.DraftItemRes{
			Title:      it.Title,
			Price:      it.Price,
			CategoryID: it.CategoryID,
		}

		if ok {
			item.CategoryIcon = cat.Icon
			item.CategoryColor = cat.ColorCode
			item.CategoryName = cat.CategoryName
		}

		items = append(items, item)
	}

	return m.DraftRes{
		DraftID:   d.DraftID,
		Status:    d.Status,
		Title:     d.Title,
		Date:      d.Date,
		Items:     items,
		Error:     d.Error,
		CreatedAt: d.CreatedAt,
	}
}

func UniqueCategoryIDsFromDrafts(drafts []m.DraftTxn) []string {
	set := make(map[string]struct{})

	for _, d := range drafts {
		for _, item := range d.Items {
			if item.CategoryID == "" {
				continue
			}
			set[item.CategoryID] = struct{}{}
		}
	}

	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}

	return ids
}
