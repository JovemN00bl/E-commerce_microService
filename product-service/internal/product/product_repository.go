package product

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetById(ctx context.Context, id int) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
}

var ErrProductNotFound = errors.New("Produto não encontrado")

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(dbPool *pgxpool.Pool) Repository {
	return &postgresRepository{db: dbPool}
}

func (r *postgresRepository) Create(ctx context.Context, product *Product) error {

	query := `
	INSERT INTO products (name, description, price, stock_quantity, created_at )
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
`

	err := r.db.QueryRow(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.CreateAt,
	).Scan(&product.ID)

	if err != nil {
		log.Printf("Erro ao criar produto no banco: %v", err)
		return err
	}

	return nil
}

func (r *postgresRepository) GetById(ctx context.Context, id int) (*Product, error) {
	query := `
	SELECT id, name, description, price, stock_quantity, created_at
	FROM products
	WHERE id = $1
`
	product := &Product{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CreateAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		log.Printf("Erro ao buscar produto por id: %v", err)
		return nil, err
	}

	return product, nil

}

func (r *postgresRepository) List(ctx context.Context) ([]*Product, error) {

	query := `
	SELECT id, name, description, price, stock_quantity, created_at
	FROM products
	ORDER BY name ASC 
`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Printf("Erro ao listar produtos: %v", err)
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)

	for rows.Next() {
		product := &Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.CreateAt,
		)
		if err != nil {
			log.Printf("Erro ao escanear linha do produto: %v", err)
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Erro após iteração das linhas de produtos: %v", err)
		return nil, err
	}

	return products, nil
}
