package domain

import (
	"context"
	"fmt"
)

type transactionableKey string

const txExistsKey = transactionableKey("transaction-exists")

type TransactionError struct {
	Cause error
}

func (D *TransactionError) Error() string {
	return fmt.Sprintf("transaction error: %s", D.Cause.Error())
}

func (D *TransactionError) TransactionErrCause() string {
	return D.Cause.Error()
}

type TransactionFunc func(ctx context.Context) error

type Transactioner interface {
	WithTx(context.Context, TransactionFunc) error
}

func InTX(ctx context.Context) bool {
	if b, ok := ctx.Value(txExistsKey).(bool); ok {
		return b
	}

	return false
}

func ContextWithTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, txExistsKey, true)
}
