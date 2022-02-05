package v1_test

import (
	"context"
	v12 "github.com/didopimentel/go-saga-poc/protogen/payments/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	dial, err := grpc.DialContext(ctx, ":7010", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := v12.NewPaymentsAPIClient(dial)

	response, err := client.CreatePayment(ctx, &v12.CreatePaymentRequest{OrderId: 1})
	require.NoError(t, err)
	require.Equal(t, int64(1), response.OrderId)
}
