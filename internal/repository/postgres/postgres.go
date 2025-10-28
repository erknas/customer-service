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

func (p *postgresRepository) InsertCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	const query = `INSERT INTO customers(username, fullname) VALUES ($1, $2) RETURNING id, created_at`

	result := &models.Customer{
		Username: customer.Username,
		Fullname: customer.Fullname,
	}

	if err := p.pool.QueryRow(
		ctx, query,
		customer.Username,
		customer.Fullname,
	).Scan(
		&result.ID,
		&result.CreatedAt,
	); err != nil {
		return nil, err
	}

	return result, nil
}
