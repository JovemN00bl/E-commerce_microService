package order

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"E-commerce_micro/order_service/internal/event"
)

var (
	ErrEmptyCart         = errors.New("O carrinho está vazio")
	ErrProductNotFound   = errors.New("Produto não encontrado ou serviço indisponível")
	ErrInsufficientStock = errors.New("Estoque insuficiente para um dos produtos")
)

type Service interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error)
	GetById(ctx context.Context, id string) (*Order, error)
}

type CreateOrderInput struct {
	UserId string
	Items  []OrderItemInput
}

type OrderItemInput struct {
	ProductId string
	Quantity  int
}

type service struct {
	repo           repository
	productClient  ProductClient
	eventPublisher event.Publisher
}

func NewService(repo repository, client ProductClient, eventPublisher event.Publisher) Service {
	return &service{repo: repo, productClient: client, eventPublisher: eventPublisher}
}

func (s *service) CreateOrder(ctx context.Context, input CreateOrderInput) (*Order, error) {
	if len(input.Items) == 0 {
		return nil, ErrEmptyCart
	}

	var orderItems []*OrderItem
	var totalPrice float64

	for _, itemInput := range input.Items {
		price, inStock, err := s.productClient.CheckStock(ctx, itemInput.ProductId)
		if err != nil {
			log.Printf("Erro ao consultar produto %s: %v", itemInput.ProductId, err)
			return nil, ErrProductNotFound
		}

		if !inStock {
			log.Printf("Estoque insuficiente para o produto %s: ", itemInput.ProductId)
			return nil, ErrInsufficientStock
		}

		itemTotal := price * (float64(itemInput.Quantity))
		totalPrice += itemTotal

		orderItems = append(orderItems, &OrderItem{
			ProductID:   itemInput.ProductId,
			Quantity:    itemInput.Quantity,
			PriceAtTime: price,
		})

	}

	order := NewOrder(input.UserId, orderItems, totalPrice)

	err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	go s.publishOrderCreated(order)

	return order, nil
}

func (s *service) GetById(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetById(ctx, id)
}

func (s *service) publishOrderCreated(order *Order) {

	payload := map[string]interface{}{
		"order_id":    order.ID,
		"user_id":     order.UserID,
		"total_price": order.TotalPrice,
		"status":      order.Status,
		"created_at":  order.CreatedAt,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao criar JSON do evento: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.eventPublisher.Publish(ctx, "order_created", jsonBody)
	if err != nil {
		log.Printf("ERRO CRÍTICO: falha ao publicar evento de pedido criado: %v ", err)
	}

}
