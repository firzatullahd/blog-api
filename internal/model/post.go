package model

import "time"

type Post struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status"`
}

type FilterFindPost struct {
	ID uint64 `db:"id"`
}

type FilterFindRPost struct {
	PostID uint64 `db:"post_id"`
}

type PostResult struct {
	ID          uint64     `json:"ID"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	PublishDate *time.Time `json:"publishDate"`
	Tags        []string   `json:"tags"`
}
