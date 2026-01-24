package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/boykush/foresee/services/transaction/internal/domain"
	"github.com/boykush/foresee/services/transaction/internal/infra/storage"
)

// CSV column indices
const (
	colID       = 0
	colDate     = 1
	colTitle    = 2
	colAmount   = 3
	colCategory = 4
)

// categoryMap maps Japanese category names to domain.Category
var categoryMap = map[string]domain.Category{
	"収入":     domain.CategoryIncome,
	"食費":     domain.CategoryFood,
	"日用品":    domain.CategoryDailyNecessities,
	"趣味・娯楽":  domain.CategoryHobby,
	"交際費":    domain.CategorySocial,
	"交通費":    domain.CategoryTransportation,
	"衣服・美容":  domain.CategoryClothing,
	"健康・医療":  domain.CategoryHealth,
	"自動車":    domain.CategoryAutomobile,
	"教養・教育":  domain.CategoryEducation,
	"特別な支出":  domain.CategorySpecial,
	"現金・カード": domain.CategoryCashCard,
	"水道・光熱費": domain.CategoryUtilities,
	"通信費":    domain.CategoryCommunication,
	"住宅":     domain.CategoryHousing,
	"税・社会保障": domain.CategoryTaxSocialSecurity,
	"保険":     domain.CategoryInsurance,
	"その他":    domain.CategoryOther,
}

// TransactionRepository implements domain.TransactionRepository using S3 storage
type TransactionRepository struct {
	storage *storage.Client
	csvKey  string
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(storage *storage.Client, csvKey string) *TransactionRepository {
	return &TransactionRepository{
		storage: storage,
		csvKey:  csvKey,
	}
}

// List returns all transactions from the CSV file
func (r *TransactionRepository) List(ctx context.Context) ([]domain.Transaction, error) {
	body, err := r.storage.GetObject(ctx, r.csvKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get CSV from storage: %w", err)
	}
	defer body.Close()

	return r.parseCSV(body)
}

func (r *TransactionRepository) parseCSV(reader io.Reader) ([]domain.Transaction, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return []domain.Transaction{}, nil
	}

	// Skip header row
	transactions := make([]domain.Transaction, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		t, err := r.parseRow(records[i])
		if err != nil {
			// Log error and continue with next row
			continue
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) parseRow(row []string) (domain.Transaction, error) {
	if len(row) < 5 {
		return domain.Transaction{}, fmt.Errorf("invalid row length: expected 5, got %d", len(row))
	}

	id := row[colID]
	title := row[colTitle]

	date, err := parseDate(row[colDate])
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("invalid date: %w", err)
	}

	amount, err := parseAmount(row[colAmount])
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("invalid amount: %w", err)
	}

	category := parseCategory(row[colCategory])

	return domain.NewTransaction(id, date, title, amount, category)
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(s))
}

func parseAmount(s string) (domain.Money, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")

	// Check if there's a decimal point
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		// Integer only
		units, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return domain.Money{}, err
		}
		return domain.NewMoney(units, 0)
	}

	// Has decimal part
	units, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return domain.Money{}, err
	}

	// Parse decimal part and convert to nanos
	decimalStr := parts[1]
	// Pad or truncate to 9 digits
	if len(decimalStr) > 9 {
		decimalStr = decimalStr[:9]
	}
	for len(decimalStr) < 9 {
		decimalStr += "0"
	}

	nanos, err := strconv.ParseInt(decimalStr, 10, 32)
	if err != nil {
		return domain.Money{}, err
	}

	// If units is negative, nanos should also be negative
	if units < 0 || (units == 0 && strings.HasPrefix(s, "-")) {
		nanos = -nanos
	}

	return domain.NewMoney(units, int32(nanos))
}

func parseCategory(s string) domain.Category {
	s = strings.TrimSpace(s)
	if cat, ok := categoryMap[s]; ok {
		return cat
	}
	return domain.CategoryUnspecified
}
