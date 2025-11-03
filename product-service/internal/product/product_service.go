package product

import "context"

type Service interface {
	Create(ctx context.Context, price float64, stock int) (*Product, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
