package domain

import "context"

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// List returns all transactions
	List(ctx context.Context) ([]Transaction, error)
}
