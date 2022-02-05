package persistence

import (
	"context"
	"fmt"
	"github.com/didopimentel/go-saga-poc/domain"
	"github.com/didopimentel/go-saga-poc/domain/entities"
)

type Payments struct {
	domain.Transactioner
	Q querier
}

const paymentsArray = "id, order_id"

func scanPayment(scanner scanner) (entities.Payment, error) {
	order := entities.Payment{}

	err := scanner.Scan(&order.ID, &order.OrderID)

	return order, err
}

func (e *Payments) CreatePayment(ctx context.Context, orderID int64) (entities.Payment, error) {
	query := fmt.Sprintf("INSERT INTO payments (order_id) VALUES ($1) RETURNING %s", paymentsArray)

	payment, err := scanPayment(e.Q.QueryRow(ctx, query, orderID))

	if err != nil {
		return entities.Payment{}, err
	}

	return payment, nil
}

func (e *Payments) DeletePayment(ctx context.Context, paymentID int64) error {
	query := "DELETE FROM payments WHERE id = $1"

	_, err := e.Q.Exec(ctx, query, paymentID)

	return err
}
