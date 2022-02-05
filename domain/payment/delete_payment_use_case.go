package payment

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type DeletePaymentUseCasePersistenceGateway interface {
	DeletePayment(context.Context, int64) error
}

type DeletePaymentUseCase struct {
	persistenceGateway DeletePaymentUseCasePersistenceGateway
}

func NewDeletePaymentUseCase(persistenceGateway DeletePaymentUseCasePersistenceGateway) *DeletePaymentUseCase {
	return &DeletePaymentUseCase{persistenceGateway: persistenceGateway}
}

type DeletePaymentInput struct {
	PaymentID int64
}
type DeletePaymentOutput struct {
	Payment entities.Delivery
}

func (u *DeletePaymentUseCase) DeletePayment(ctx context.Context, input DeletePaymentInput) (DeletePaymentOutput, error) {
	err := u.persistenceGateway.DeletePayment(ctx, input.PaymentID)

	return DeletePaymentOutput{}, err
}
