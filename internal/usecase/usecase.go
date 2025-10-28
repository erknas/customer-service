package usecase

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
)

type UserRepository interface {
	InsertCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)
}

type customerUseCase struct {
	repo UserRepository
}

func New(repo UserRepository) *customerUseCase {
	return &customerUseCase{
		repo: repo,
	}
}
