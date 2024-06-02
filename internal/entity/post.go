package entity

import (
	"time"
)

type Post struct {
	ID      uint64 `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	//Tags    []uint64 `db:"tags"`
	Status string `db:"status"`

	PublishDate *time.Time `db:"publish_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type PostStatus int

const (
	StatusDraft PostStatus = iota + 1
	StatusPublish
)

func (e PostStatus) String() string {
	switch e {
	case StatusDraft:
		return "draft"
	case StatusPublish:
		return "publish"
	default:
		return "unknown"
	}
}
