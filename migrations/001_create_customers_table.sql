-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_name VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(100),
    city VARCHAR(100),
    birth_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    CONSTRAINT valid_birth_date CHECK (
        birth_date <= CURRENT_DATE
        AND birth_date >= '1900-01-01'
    )
);

CREATE INDEX idx_customers_created_at ON customers (created_at);

CREATE OR REPLACE FUNCTION update_updated_at_column () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_customers_updated_at BEFORE
UPDATE ON customers FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

-- +goose StatementEnd