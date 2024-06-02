package entity

import "time"

type RPost struct {
	ID        uint64     `db:"id"`
	PostID    uint64     `db:"post_id"`
	TagID     uint64     `db:"tag_id"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
