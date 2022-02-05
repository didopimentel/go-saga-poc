package persistence

import (
	"context"
)

type Health struct {
	Q querier
}

// Check does a simple query over the DB connection to validate its health
func (h *Health) Check(ctx context.Context) error {
	_, err := h.Q.Exec(ctx, "SELECT 1;")
	return err
}
