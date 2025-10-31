package usecase

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
	"go.uber.org/zap"
)

type CustomerRepository interface {
	Insert(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	Customer(ctx context.Context, id string) (*models.Customer, error)
	Update(ctx context.Context, id string, customer *models.UpdatedCustomer) (*models.Customer, error)
}

type customerUseCase struct {
	repo   CustomerRepository
	logger *zap.Logger
}

func New(repo CustomerRepository, logger *zap.Logger) *customerUseCase {
	return &customerUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (c *customerUseCase) AddCustomer(ctx context.Context, userName, fullName, city, birthDateStr string) (*models.Customer, error) {
	customer, err := models.NewCustomer(userName, fullName, city, birthDateStr)
	if err != nil {
		return nil, err
	}

	result, err := c.repo.Insert(ctx, customer)
	if err != nil {
		c.logger.Error("failed to insert customer", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (c *customerUseCase) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	customer, err := c.repo.Customer(ctx, id)
	if err != nil {
		c.logger.Error("failed to get customer", zap.Error(err))
		return nil, err
	}

	return customer, nil
}

func (c *customerUseCase) UpdateCustomer(ctx context.Context, id string, userName, fullName, city, birthDateStr *string) (*models.Customer, error) {
	customer, err := models.NewUpdatedCustomer(userName, fullName, city, birthDateStr)
	if err != nil {
		return nil, err
	}

	result, err := c.repo.Update(ctx, id, customer)
	if err != nil {
		c.logger.Error("failed to udpate customer", zap.Error(err))
		return nil, err
	}

	return result, nil
}
