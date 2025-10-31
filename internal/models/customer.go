package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const layout = "2006-01-02"

type Customer struct {
	ID        uuid.UUID
	Username  string
	Fullname  string
	City      string
	BirthDate time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdatedCustomer struct {
	Username  *string
	Fullname  *string
	City      *string
	BirthDate *time.Time
}

func NewCustomer(userName, fullName, city, birthDateStr string) (*Customer, error) {
	birthDate, err := time.Parse(layout, birthDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	return &Customer{
		Username:  userName,
		Fullname:  fullName,
		City:      city,
		BirthDate: birthDate,
	}, nil
}

func NewUpdatedCustomer(userName, fullName, city, birthDateStr *string) (*UpdatedCustomer, error) {
	if birthDateStr == nil {
		return &UpdatedCustomer{
			Username:  userName,
			Fullname:  fullName,
			City:      city,
			BirthDate: nil,
		}, nil
	}

	birthDate, err := time.Parse(layout, *birthDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	return &UpdatedCustomer{
		Username:  userName,
		Fullname:  fullName,
		City:      city,
		BirthDate: &birthDate,
	}, nil
}
