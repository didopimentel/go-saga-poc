package payments

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/entities"
	v1 "github.com/didopimentel/go-saga-poc/protogen/payments/api/v1"
)

type Gateway struct {
	cli v1.PaymentsAPIClient
}

func NewGateway(client v1.PaymentsAPIClient) *Gateway {
	return &Gateway{cli: client}
}

func (g *Gateway) CreatePayment(ctx context.Context, orderID int64) (entities.Payment, error) {
	response, err := g.cli.CreatePayment(ctx, &v1.CreatePaymentRequest{OrderId: orderID})
	if err != nil {
		return entities.Payment{}, err
	}

	return entities.Payment{
		ID:      response.Id,
		OrderID: response.OrderId,
	}, nil
}

func (g *Gateway) DeletePayment(ctx context.Context, paymentID int64) error {
	_, err := g.cli.DeletePayment(ctx, &v1.DeletePaymentRequest{Id: paymentID})

	return err
}
