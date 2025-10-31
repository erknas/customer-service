package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/erknas/customer-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueConstraintCode = "23505"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		pool: pool,
	}
}

func (p *postgresRepository) Insert(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	const query = `
	INSERT INTO customers(user_name, full_name, city, birth_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id, is_active, created_at, updated_at
	`

	result := &models.Customer{
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
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueConstraintCode {
				return nil, fmt.Errorf("user already exists")
			}
		}
		return nil, err
	}

	return result, nil
}

func (p *postgresRepository) Customer(ctx context.Context, id string) (*models.Customer, error) {
	const query = `
	SELECT id, user_name, full_name, city, birth_date, is_active, created_at, updated_at
	FROM customers
	WHERE id = $1
	`

	var result models.Customer
	if err := p.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Fullname,
		&result.City,
		&result.BirthDate,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &result, nil
}

func (p *postgresRepository) Update(ctx context.Context, id string, customer *models.UpdatedCustomer) (*models.Customer, error) {
	columns := []string{}
	args := []any{}
	argPos := 1

	if customer.Username != nil {
		columns = append(columns, fmt.Sprintf("user_name = $%d", argPos))
		args = append(args, customer.Username)
		argPos++
	}

	if customer.Fullname != nil {
		columns = append(columns, fmt.Sprintf("full_name = $%d", argPos))
		args = append(args, customer.Fullname)
		argPos++
	}

	if customer.City != nil {
		columns = append(columns, fmt.Sprintf("city = $%d", argPos))
		args = append(args, customer.City)
		argPos++
	}

	if customer.BirthDate != nil {
		columns = append(columns, fmt.Sprintf("birth_date = $%d", argPos))
		args = append(args, customer.BirthDate)
		argPos++
	}

	query := `UPDATE customers SET ` + strings.Join(columns, ", ") + fmt.Sprintf(" WHERE id = $%d ", argPos) + "RETURNING id, user_name, full_name, city, birth_date, is_active, created_at, updated_at"
	args = append(args, id)

	var result models.Customer
	if err := p.pool.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Fullname,
		&result.City,
		&result.BirthDate,
		&result.IsActive,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueConstraintCode {
				return nil, fmt.Errorf("user already exists")
			}
		}
		return nil, err
	}

	return &result, nil
}
