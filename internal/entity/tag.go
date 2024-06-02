package entity

import (
	"time"
)

type Tag struct {
	ID    uint64 `db:"id"`
	Label string `db:"label"`
	//Posts []uint64 `db:"posts"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
