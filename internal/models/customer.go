package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID
	Username  string
	Fullname  string
	CreatedAt time.Time
}
