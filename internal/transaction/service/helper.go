package transactionsvc

import (
	"errors"
	"fmt"
	"strings"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	draft "github.com/cp25sy5-modjot/main-service/internal/draft"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"github.com/google/uuid"
)

func ReplaceTransactionItems(
	txID string,
	input []m.TransactionItemInput,
) ([]e.TransactionItem, error) {

	if len(input) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	items := make([]e.TransactionItem, 0, len(input))

	for _, it := range input {
		items = append(items, e.TransactionItem{
			TransactionID: txID,
			ItemID:        uuid.New().String(),
			Title:         it.Title,
			Price:         it.Price,
			CategoryID:    it.CategoryID,
		})
	}

	return items, nil
}

func mapToDraft(
	resp *pb.TransactionResponseV2,
	categories []e.Category,
	userID string,
) (*draft.DraftTxn, error) {

	if len(resp.Items) == 0 {
		return nil, errors.New("no transaction items")
	}

	var items []draft.DraftItem

	for _, res := range resp.Items {

		match := matchCategoryFromName(categories, res.Category)
		if match == nil {
			return nil, fmt.Errorf("category not found: %s", res.Category)
		}

		items = append(items, draft.DraftItem{
			Title:      res.Title,
			Price:      res.Price,
			CategoryID: match.CategoryID,
		})
	}

	date, err := ParseAIDate(resp.Date)
	if err != nil {
		return nil, err
	}

	return &draft.DraftTxn{
		UserID:    userID,
		Status:    draft.DraftStatusWaitingConfirm,
		Title:     resp.Title,
		Date:      &date,
		Items:     items,
		CreatedAt: time.Now(),
	}, nil
}

var bkkLoc = loadBKK()

func loadBKK() *time.Location {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// fallback สำหรับ container ที่ไม่มี tzdata
		return time.FixedZone("Asia/Bangkok", 7*60*60)
	}
	return loc
}

func ParseAIDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return time.Time{}, errors.New("empty date string")
	}

	// 🧹 AI ชอบส่ง timezone มั่ว → ตัดทิ้งก่อน
	s = stripTimezone(s)

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {

		var t time.Time
		var err error

		if layout == time.RFC3339 {
			t, err = time.Parse(layout, s)
			if err == nil {
				return t.UTC(), nil
			}
			continue
		}

		t, err = time.ParseInLocation(layout, s, bkkLoc)
		if err == nil {

			// date only → set noon
			if layout == "2006-01-02" {
				t = time.Date(
					t.Year(),
					t.Month(),
					t.Day(),
					12, 0, 0, 0,
					bkkLoc,
				)
			}

			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format from AI: %q", s)
}

func stripTimezone(s string) string {
	s = strings.TrimSuffix(s, "Z")

	suffixes := []string{
		"+0000", "+00:00",
		"+0700", "+07:00",
	}

	for _, suf := range suffixes {
		s = strings.TrimSuffix(s, suf)
	}

	return s
}

func (s *service) validateUpdateItems(
	userID string,
	items []m.TransactionItemInput,
) error {

	if len(items) == 0 {
		return errors.New("at least one item is required")
	}

	// ดึง category ของ user
	categories, err := s.catrepo.FindAllByUserID(userID)
	if err != nil {
		return err
	}

	// map categoryID -> true
	categoryMap := map[string]bool{}
	for _, c := range categories {
		categoryMap[c.CategoryID] = true
	}

	for _, it := range items {

		if it.Title == "" {
			return errors.New("item title is required")
		}

		if it.Price < 0 {
			return errors.New("item price must be positive")
		}

		if !categoryMap[it.CategoryID] {
			return errors.New("invalid category")
		}
	}

	return nil
}
