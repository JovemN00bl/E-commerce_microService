package product

import "time"

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"Description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stockQuantity"`
	CreateAt      time.Time `json:"createAt"`
}

func NewProduct(name, description string, price float64, stockQuantity int) *Product {
	return &Product{Name: name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		CreateAt:      time.Now().UTC()}
}
