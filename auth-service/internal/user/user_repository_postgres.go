package user

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type postgresRepositoy struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(dbPool *pgxpool.Pool) Repository {
	return &postgresRepositoy{db: dbPool}
}

func (r *postgresRepositoy) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (email, password_hash, created_at)
	VALUES ($1, $2, $3)
	RETURNING id 
`

	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)

	if err != nil {
		log.Println("Erro ao criar usuario no banco: %v", err)
		return err
	}
	return nil
}

func (r *postgresRepositoy) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE email = $1
`

	user := &User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		log.Printf("Erro ao buscar usuário por email: %v", err)
		return nil, err

	}
	return user, nil

}

var ErrUserNotFound = errors.New("usuário não encontrado")
