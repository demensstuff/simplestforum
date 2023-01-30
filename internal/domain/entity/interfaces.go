package entity

import (
	"context"
)

// Transactioner is anything which can spawn a new transaction (if any), i.e. storage connection.
type Transactioner interface {
	NewTransaction(context.Context) (AbstractTransaction, error)
}

// Transactionable is something which abstracts a transaction to propagate it, i.e. service.
type Transactionable interface {
	DoTransaction(Session, func() error) error
}
