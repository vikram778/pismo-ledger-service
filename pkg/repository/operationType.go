package repository

import (
	"context"
	"go.uber.org/zap"
	"pismo-ledger-service/pkg/log"
	"pismo-ledger-service/pkg/models"
)

// GetOperationType gets the operation type details based on the operation ID
// Param - context object
// Param - id type int64 for which operation details needs to be fetched
// Returns - operation model
// Returns - error if any
func (r *dBOps) GetOperationType(ctx context.Context, id int64) (*models.Operations, error) {

	var operationType models.Operations

	if err := r.db.Get(&operationType, getOperationTypeQuery, id); err != nil {
		log.Info("fetch operation type error", zap.Error(err))
		return nil, err
	}
	return &operationType, nil
}
