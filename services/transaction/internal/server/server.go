package server

import (
	"context"

	transactionv1 "github.com/boykush/foresee/services/transaction/gen/go"
	"github.com/boykush/foresee/services/transaction/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the TransactionServiceServer interface
type Server struct {
	transactionv1.UnimplementedTransactionServiceServer
	repo domain.TransactionRepository
}

// NewServer creates a new instance of the transaction service server
func NewServer(repo domain.TransactionRepository) *Server {
	return &Server{
		repo: repo,
	}
}

// HealthCheck implements the health check endpoint
func (s *Server) HealthCheck(ctx context.Context, req *transactionv1.HealthCheckRequest) (*transactionv1.HealthCheckResponse, error) {
	return &transactionv1.HealthCheckResponse{
		Status: transactionv1.HealthCheckResponse_SERVING,
	}, nil
}

// ListTransactions returns all transactions from storage
func (s *Server) ListTransactions(ctx context.Context, req *transactionv1.ListTransactionsRequest) (*transactionv1.ListTransactionsResponse, error) {
	transactions, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	pbTransactions := make([]*transactionv1.Transaction, len(transactions))
	for i, t := range transactions {
		pbTransactions[i] = toProtoTransaction(t)
	}

	return &transactionv1.ListTransactionsResponse{
		Transactions: pbTransactions,
	}, nil
}

func toProtoTransaction(t domain.Transaction) *transactionv1.Transaction {
	return &transactionv1.Transaction{
		Id:    t.ID(),
		Date:  timestamppb.New(t.Date()),
		Title: t.Title(),
		Amount: &transactionv1.Money{
			Units: t.Amount().Units(),
			Nanos: t.Amount().Nanos(),
		},
		Category: toProtoCategory(t.Category()),
	}
}

func toProtoCategory(c domain.Category) transactionv1.Category {
	switch c {
	case domain.CategoryIncome:
		return transactionv1.Category_CATEGORY_INCOME
	case domain.CategoryFood:
		return transactionv1.Category_CATEGORY_FOOD
	case domain.CategoryDailyNecessities:
		return transactionv1.Category_CATEGORY_DAILY_NECESSITIES
	case domain.CategoryHobby:
		return transactionv1.Category_CATEGORY_HOBBY
	case domain.CategorySocial:
		return transactionv1.Category_CATEGORY_SOCIAL
	case domain.CategoryTransportation:
		return transactionv1.Category_CATEGORY_TRANSPORTATION
	case domain.CategoryClothing:
		return transactionv1.Category_CATEGORY_CLOTHING
	case domain.CategoryHealth:
		return transactionv1.Category_CATEGORY_HEALTH
	case domain.CategoryAutomobile:
		return transactionv1.Category_CATEGORY_AUTOMOBILE
	case domain.CategoryEducation:
		return transactionv1.Category_CATEGORY_EDUCATION
	case domain.CategorySpecial:
		return transactionv1.Category_CATEGORY_SPECIAL
	case domain.CategoryCashCard:
		return transactionv1.Category_CATEGORY_CASH_CARD
	case domain.CategoryUtilities:
		return transactionv1.Category_CATEGORY_UTILITIES
	case domain.CategoryCommunication:
		return transactionv1.Category_CATEGORY_COMMUNICATION
	case domain.CategoryHousing:
		return transactionv1.Category_CATEGORY_HOUSING
	case domain.CategoryTaxSocialSecurity:
		return transactionv1.Category_CATEGORY_TAX_SOCIAL_SECURITY
	case domain.CategoryInsurance:
		return transactionv1.Category_CATEGORY_INSURANCE
	case domain.CategoryOther:
		return transactionv1.Category_CATEGORY_OTHER
	default:
		return transactionv1.Category_CATEGORY_UNSPECIFIED
	}
}
