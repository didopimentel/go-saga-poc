package delivery

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type CreateDeliveryUseCasePersistenceGateway interface {
	CreateDelivery(context.Context, int64) (entities.Delivery, error)
}

type CreateDeliveryUseCase struct {
	persistenceGateway CreateDeliveryUseCasePersistenceGateway
}

func NewCreateDeliveryUseCase(persistenceGateway CreateDeliveryUseCasePersistenceGateway) *CreateDeliveryUseCase {
	return &CreateDeliveryUseCase{persistenceGateway: persistenceGateway}
}

type CreateDeliveryInput struct {
	OrderID int64
}
type CreateDeliveryOutput struct {
	Delivery entities.Delivery
}

func (u *CreateDeliveryUseCase) CreateDelivery(ctx context.Context, input CreateDeliveryInput) (CreateDeliveryOutput, error) {
	delivery, err := u.persistenceGateway.CreateDelivery(ctx, input.OrderID)
	if err != nil {
		return CreateDeliveryOutput{}, err
	}
	return CreateDeliveryOutput{
		Delivery: delivery,
	}, nil
}
