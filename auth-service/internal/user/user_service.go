package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrEmailAlreadyExists = errors.New("e-mail já cadastrado")
)

type Service interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (TokenString string, err error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{
		jwtSecret: jwtSecret,
		repo:      repo}
}

func (s *service) Register(ctx context.Context, email, password string) error {
	//todo may add more validations, like email is valid ? or strong password

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := NewUser(email, string(hash))

	err = s.repo.Create(ctx, user)
	if err != nil {
		//todo check if the error is a duplicate email, and if so, return EmailAlreadyExists
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, email, password string) (TokenString string, err error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateJwt(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) generateJwt(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
