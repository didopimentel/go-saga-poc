package v1

import (
	"context"
	"github.com/didopimentel/go-saga-poc/gateways/persistence"
	v12 "github.com/didopimentel/go-saga-poc/protogen/payments/api/v1"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

var _ v12.PaymentsAPIServer = &API{}

// API implements emv1pb.EventsManagerAPIServer
type API struct {
	*Repository
	*PaymentsAPI
}

type Repository struct {
	Payments *persistence.Payments
	Health   *persistence.Health
}

// GetHealth lets clients know if Events Manager Server is healthy to respond requests
func (a *API) GetHealth(ctx context.Context, request *v12.GetHealthRequest) (*v12.GetHealthResponse, error) {
	// TODO: maybe we don't want to be unhealthy if PG is down! Let's start like this though
	if err := a.Repository.Health.Check(ctx); err != nil {
		return nil, err
	}

	return &v12.GetHealthResponse{}, nil
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
