package order

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository interface {
	GetById(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id string, status OrderStatus) error
}

var ErrOrderNotFound = errors.New("Pedido não encontrado!")

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, order *Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryOrder := `
		INSERT INTO orders (user_id, total_price, status, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(ctx, queryOrder,
		order.UserID,
		order.TotalPrice,
		order.Status,
		order.CreatedAt,
	).Scan(&order.ID)

	if err != nil {
		log.Printf("Erro ao inserir pedido: %v", err)
		return err
	}

	queryItem := `
		INSERT INTO order_items (order_id, product_id, quantity, price_at_time)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		item.OrderID = order.ID

		_, err := tx.Exec(ctx, queryItem,
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.PriceAtTime,
		)
		if err != nil {
			log.Printf("Erro ao inserir item do pedido: %v", err)
			return err
		}
	}
	return tx.Commit(ctx)

}

func (r *postgresRepository) GetById(ctx context.Context, id string) (*Order, error) {
	queryOrder := `SELECT id, user_id, total_price, status, created_at FROM orders WHERE id = $1`

	order := &Order{}
	err := r.db.QueryRow(ctx, queryOrder, id).Scan(
		&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	queryItems := `SELECT id, order_id, product_id, quantity, price_at_time FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(ctx, queryItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &OrderItem{}
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PriceAtTime); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id string, status OrderStatus) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		log.Printf("Erro ao atualizar status do pedido %s: %v", id, err)
		return err
	}
	return nil

}
