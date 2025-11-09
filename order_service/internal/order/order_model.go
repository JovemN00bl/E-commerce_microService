package order

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID         string       `json:"ID"`
	UserID     string       `json:"user_id"`
	TotalPrice float64      `json:"total_price"`
	Status     OrderStatus  `json:"status"`
	Items      []*OrderItem `json:"items"`
	CreatedAt  time.Time    `json:"Created_At"`
}

type OrderItem struct {
	ID          string  `json:"ID"`
	OrderID     string  `json:"Order_ID"`
	ProductID   string  `json:"Product_ID"`
	Quantity    int     `json:"quantity"`
	PriceAtTime float64 `json:"price_at_time"`
}

func NewOrder(UserID string, items []*OrderItem, TotalPrice float64) *Order {
	return &Order{UserID: UserID,
		Items:      items,
		TotalPrice: TotalPrice,
		Status:     StatusPending,
		CreatedAt:  time.Now().UTC()}
}
