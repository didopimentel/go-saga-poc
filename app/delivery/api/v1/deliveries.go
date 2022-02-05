package v1

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/delivery"
	v12 "github.com/didopimentel/go-saga-poc/protogen/delivery/api/v1"
)

type DeliveryAPI struct {
	DeliveryAPIUseCases
}

func NewDeliveryAPI(uc DeliveryAPIUseCases) *DeliveryAPI {
	return &DeliveryAPI{
		DeliveryAPIUseCases: uc,
	}
}

type DeliveryAPIUseCases interface {
	CreateDelivery(ctx context.Context, input delivery.CreateDeliveryInput) (delivery.CreateDeliveryOutput, error)
}

func (a *DeliveryAPI) CreateDelivery(ctx context.Context, req *v12.CreateDeliveryRequest) (*v12.CreateDeliveryResponse, error) {
	o, err := a.DeliveryAPIUseCases.CreateDelivery(ctx, delivery.CreateDeliveryInput{
		OrderID: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return &v12.CreateDeliveryResponse{
		Id:      o.Delivery.ID,
		OrderId: o.Delivery.OrderID,
	}, nil
}
