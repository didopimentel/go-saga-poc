package order

import (
	"context"
	"errors"
	"github.com/didopimentel/go-saga-poc/domain"
	"github.com/didopimentel/go-saga-poc/domain/entities"
	"github.com/didopimentel/go-saga-poc/extensions/saga"
	"log"
)

type CreateOrderUseCasePersistenceGateway interface {
	CreateOrder(context.Context, int64) (entities.Order, error)
}

type CreateOrderUseCasePaymentGateway interface {
	CreatePayment(ctx context.Context, orderID int64) (entities.Payment, error)
	DeletePayment(ctx context.Context, paymentID int64) error
}

type CreateOrderUseCaseDeliveriesGateway interface {
	CreateDelivery(ctx context.Context, orderID int64) (entities.Delivery, error)
}

type CreateOrderUseCase struct {
	persistenceGateway CreateOrderUseCasePersistenceGateway
	paymentsGateway    CreateOrderUseCasePaymentGateway
	deliveriesGateway  CreateOrderUseCaseDeliveriesGateway
	tx                 domain.Transactioner
}

func NewCreateOrderUseCase(persistenceGateway CreateOrderUseCasePersistenceGateway,
	tx domain.Transactioner,
	paymentsGateway CreateOrderUseCasePaymentGateway,
	deliveriesGateway CreateOrderUseCaseDeliveriesGateway) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		persistenceGateway: persistenceGateway,
		tx:                 tx,
		paymentsGateway:    paymentsGateway,
		deliveriesGateway:  deliveriesGateway,
	}
}

type CreateOrderInput struct {
	Amount int64
}
type CreateOrderOutput struct {
	Order entities.Order
}

func (u *CreateOrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	output := CreateOrderOutput{}
	err := u.tx.WithTx(ctx, func(ctx context.Context) error {

		coordinator := u.getCreateOrderSagaCoordinator()
		
		sagaContext := context.WithValue(ctx, saga.ParamKey, input)
		result, ok := coordinator.Execute(sagaContext)
		if !ok {
			for _, e := range coordinator.GetErrors() {
				log.Println(e.Error())
			}
			return errors.New("could not create order")
		}

		output.Order = result.(entities.Order)
		return nil
	})
	if err != nil {
		return CreateOrderOutput{}, err
	}

	return output, nil
}

func (u *CreateOrderUseCase) getCreateOrderSagaCoordinator() *saga.Coordinator {
	steps := []saga.Step{
		{
			Command: func(ctx context.Context) (interface{}, error) {
				reqInput := ctx.Value(saga.ParamKey).(CreateOrderInput)
				createdOrder, err := u.persistenceGateway.CreateOrder(ctx, reqInput.Amount)
				if err != nil {
					return nil, err
				}

				return createdOrder, nil
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				return nil, nil
			},
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				reqInput := ctx.Value(saga.ParamKey).(entities.Order)
				payment, err := u.paymentsGateway.CreatePayment(ctx, reqInput.ID)
				if err != nil {
					return nil, err
				}

				reqInput.PaymentID = payment.ID
				return reqInput, nil
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				reqInput := ctx.Value(saga.ParamKey).(entities.Order)
				err := u.paymentsGateway.DeletePayment(ctx, reqInput.PaymentID)
				if err != nil {
					return nil, err
				}
				return nil, nil
			},
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				reqInput := ctx.Value(saga.ParamKey).(entities.Order)
				delivery, err := u.deliveriesGateway.CreateDelivery(ctx, reqInput.ID)
				if err != nil {
					return nil, err
				}

				reqInput.DeliveryID = delivery.ID
				return reqInput, nil
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				return nil, nil
			},
		},
	}
	createOrderSaga := saga.NewSaga(steps)
	return saga.NewCoordinator(createOrderSaga)
}
