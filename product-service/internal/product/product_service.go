package product

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, name, description string, price float64, stock int) (*Product, error)
	GetById(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
}

var (
	ErrProductNameRequired = errors.New("o nome do produto é obrigatorio.")
	ErrInvalidPrice        = errors.New("o preço deve ser positivo.")
	ErrInvalidStock        = errors.New("O estoque não pode ser negativo.")
)

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, name, description string, price float64, stock int) (*Product, error) {
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if price < 0 {
		return nil, ErrInvalidPrice
	}
	if stock <= 0 {
		return nil, ErrInvalidStock
	}

	product := NewProduct(name, description, price, stock)

	err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *service) GetById(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetById(ctx, id)
}

func (s *service) List(ctx context.Context) ([]*Product, error) {
	return s.repo.List(ctx)
}
