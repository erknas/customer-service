package usecase

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
	"go.uber.org/zap"
)

type customerRepository interface {
	InsertCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
}

type customerUseCase struct {
	repo   customerRepository
	logger *zap.Logger
}

func New(repo customerRepository, logger *zap.Logger) *customerUseCase {
	return &customerUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (c *customerUseCase) RegisterCustomer(ctx context.Context, userName, fullName, city, birthDateStr string) (models.Customer, error) {
	customer, err := models.NewCustomer(userName, fullName, city, birthDateStr)
	if err != nil {
		return models.Customer{}, err
	}

	result, err := c.repo.InsertCustomer(ctx, customer)
	if err != nil {
		c.logger.Error("failed to insert customer", zap.Error(err))
		return models.Customer{}, err
	}

	return result, nil
}
