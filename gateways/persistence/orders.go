package persistence

import (
	"context"
	"fmt"
	"github.com/didopimentel/go-saga-poc/domain"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type Orders struct {
	domain.Transactioner
	Q querier
}

const ordersArray = "id, amount"

func scanActivityLog(scanner scanner) (entities.Order, error) {
	order := entities.Order{}

	err := scanner.Scan(&order.ID, &order.Amount)

	return order, err
}

func (e *Orders) CreateOrder(ctx context.Context, amount int64) (entities.Order, error) {
	query := fmt.Sprintf("INSERT INTO orders (amount) VALUES ($1) RETURNING %s", ordersArray)

	order, err := scanActivityLog(e.Q.QueryRow(ctx, query, amount))

	if err != nil {
		return entities.Order{}, err
	}

	return order, nil
}
