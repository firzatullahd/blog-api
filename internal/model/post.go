package model

import "time"

type Post struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status"`
}

type FilterFindPost struct {
	IDs   []uint64 `db:"id"`
	Limit int
	Page  int
}

type FilterFindRPost struct {
	PostID uint64   `db:"post_id"`
	TagIDs []uint64 `db:"tag_id"`
}

type FilterSearchPost struct {
	TagLabel []string
	Page     int
	Limit    int
}

type PostResult struct {
	ID          uint64     `json:"ID"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	PublishDate *time.Time `json:"publishDate"`
	Tags        []string   `json:"tags"`
}

type SearchResult struct {
	Data  []PostResult `json:"data"`
	Page  int          `json:"page"`
	Count int          `json:"count"`
}
