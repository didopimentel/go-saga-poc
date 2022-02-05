package entities

type Order struct {
	ID         int64
	Amount     int64
	PaymentID  int64
	DeliveryID int64
}
