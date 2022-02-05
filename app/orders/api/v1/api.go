package v1

import (
	"context"
	"github.com/didopimentel/go-saga-poc/gateways/persistence"
	v1 "github.com/didopimentel/go-saga-poc/protogen/orders/api/v1"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

var _ v1.OrdersAPIServer = &API{}

// API implements emv1pb.EventsManagerAPIServer
type API struct {
	*Repository
	*OrdersAPI
}

type Repository struct {
	Orders *persistence.Orders
	Health *persistence.Health
}

// GetHealth lets clients know if Events Manager Server is healthy to respond requests
func (a *API) GetHealth(ctx context.Context, _ *v1.GetHealthRequest) (*v1.GetHealthResponse, error) {
	// TODO: maybe we don't want to be unhealthy if PG is down! Let's start like this though
	if err := a.Repository.Health.Check(ctx); err != nil {
		return nil, err
	}

	return &v1.GetHealthResponse{}, nil
}

func TimeToTimestamp(t *time.Time) *timestamp.Timestamp {
	if t == nil {
		return nil
	}

	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}
