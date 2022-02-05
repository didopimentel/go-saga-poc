package v1

import (
	"context"
	"github.com/didopimentel/go-saga-poc/domain/payment"
	v1 "github.com/didopimentel/go-saga-poc/protogen/payments/api/v1"
)

type PaymentsAPI struct {
	PaymentsAPIUseCases
}

func NewPaymentsAPI(uc PaymentsAPIUseCases) *PaymentsAPI {
	return &PaymentsAPI{
		PaymentsAPIUseCases: uc,
	}
}

type PaymentsAPIUseCases interface {
	CreatePayment(ctx context.Context, input payment.CreatePaymentInput) (payment.CreatePaymentOutput, error)
	DeletePayment(ctx context.Context, input payment.DeletePaymentInput) (payment.DeletePaymentOutput, error)
}

func (a *PaymentsAPI) CreatePayment(ctx context.Context, req *v1.CreatePaymentRequest) (*v1.CreatePaymentResponse, error) {
	o, err := a.PaymentsAPIUseCases.CreatePayment(ctx, payment.CreatePaymentInput{
		OrderID: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreatePaymentResponse{
		Id:      o.Payment.ID,
		OrderId: o.Payment.OrderID,
	}, nil
}

func (a *PaymentsAPI) DeletePayment(ctx context.Context, req *v1.DeletePaymentRequest) (*v1.DeletePaymentResponse, error) {
	_, err := a.PaymentsAPIUseCases.DeletePayment(ctx, payment.DeletePaymentInput{
		PaymentID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &v1.DeletePaymentResponse{}, nil
}
