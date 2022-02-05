package deliveries

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/entities"
	v1 "github.com/didopimentel/go-saga-poc/protogen/delivery/api/v1"
)

type Gateway struct {
	cli v1.DeliveryAPIClient
}

func NewGateway(client v1.DeliveryAPIClient) *Gateway {
	return &Gateway{cli: client}
}

func (g *Gateway) CreateDelivery(ctx context.Context, orderID int64) (entities.Delivery, error) {
	response, err := g.cli.CreateDelivery(ctx, &v1.CreateDeliveryRequest{OrderId: orderID})
	if err != nil {
		return entities.Delivery{}, err
	}

	return entities.Delivery{
		ID:      response.Id,
		OrderID: response.OrderId,
	}, nil
}
