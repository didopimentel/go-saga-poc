package persistence

import (
	"context"
	"fmt"
	"github.com/didopimentel/go-saga-poc/domain"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type Deliveries struct {
	domain.Transactioner
	Q querier
}

const deliveriesArray = "id, order_id"

func scanDelivery(scanner scanner) (entities.Delivery, error) {
	delivery := entities.Delivery{}

	err := scanner.Scan(&delivery.ID, &delivery.OrderID)

	return delivery, err
}

func (e *Deliveries) CreateDelivery(ctx context.Context, orderID int64) (entities.Delivery, error) {
	query := fmt.Sprintf("INSERT INTO deliveries (order_id) VALUES ($1) RETURNING %s", deliveriesArray)

	delivery, err := scanDelivery(e.Q.QueryRow(ctx, query, orderID))

	if err != nil {
		return entities.Delivery{}, err
	}

	return delivery, nil
}
