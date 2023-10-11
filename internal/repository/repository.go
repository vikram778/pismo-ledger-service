package repository

import (
	"context"
	models2 "pismo-ledger-service/pkg/models"
)

// DBOps interface defines all the operations performed on the DB
type DBOps interface {
	CreateAccount(ctx context.Context, profile *models2.Account) error
	GetAccount(ctx context.Context, id int64) (*models2.Account, error)
	GetAccountByDocument(ctx context.Context, docNo string) (*models2.Account, error)
	GetOperationType(ctx context.Context, id int64) (*models2.Operations, error)
	CreateTransaction(ctx context.Context, txn *models2.Transaction) error
}
