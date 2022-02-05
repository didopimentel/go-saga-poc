package v1_test

import (
	"context"
	v1 "github.com/didopimentel/go-saga-poc/protogen/orders/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	dial, err := grpc.DialContext(ctx, ":7000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := v1.NewOrdersAPIClient(dial)

	response, err := client.CreateOrder(ctx, &v1.CreateOrderRequest{Amount: 100})
	require.NoError(t, err)
	require.Equal(t, int64(100), response.Amount)
}
