package payment

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type CreatePaymentUseCasePersistenceGateway interface {
	CreatePayment(context.Context, int64) (entities.Payment, error)
}

type CreatePaymentUseCase struct {
	persistenceGateway CreatePaymentUseCasePersistenceGateway
}

func NewCreatePaymentUseCase(persistenceGateway CreatePaymentUseCasePersistenceGateway) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{persistenceGateway: persistenceGateway}
}

type CreatePaymentInput struct {
	OrderID int64
}
type CreatePaymentOutput struct {
	Payment entities.Payment
}

func (u *CreatePaymentUseCase) CreatePayment(ctx context.Context, input CreatePaymentInput) (CreatePaymentOutput, error) {
	payment, err := u.persistenceGateway.CreatePayment(ctx, input.OrderID)
	if err != nil {
		return CreatePaymentOutput{}, err
	}
	return CreatePaymentOutput{
		Payment: payment,
	}, nil
}
