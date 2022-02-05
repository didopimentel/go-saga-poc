package v1

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/order"
	v1 "github.com/didopimentel/go-saga-poc/protogen/orders/api/v1"
)

type OrdersAPI struct {
	OrdersAPIUseCases
}

func NewOrdersAPI(uc OrdersAPIUseCases) *OrdersAPI {
	return &OrdersAPI{
		OrdersAPIUseCases: uc,
	}
}

type OrdersAPIUseCases interface {
	CreateOrder(ctx context.Context, input order.CreateOrderInput) (order.CreateOrderOutput, error)
}

func (a *OrdersAPI) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderResponse, error) {
	o, err := a.OrdersAPIUseCases.CreateOrder(ctx, order.CreateOrderInput{
		Amount: req.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateOrderResponse{
		Id:     o.Order.ID,
		Amount: o.Order.Amount,
	}, nil
}
