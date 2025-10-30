package postgres

import (
	"context"

	"github.com/erknas/customer-service/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		pool: pool,
	}
}

func (p *postgresRepository) InsertCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	const query = `
	INSERT INTO customers(user_name, full_name, city, birth_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	result := models.Customer{
		Username:  customer.Username,
		Fullname:  customer.Fullname,
		City:      customer.City,
		BirthDate: customer.BirthDate,
	}

	if err := p.pool.QueryRow(
		ctx, query,
		customer.Username,
		customer.Fullname,
		customer.City,
		customer.BirthDate,
	).Scan(
		&result.ID,
	); err != nil {
		return models.Customer{}, err
	}

	return result, nil
}
