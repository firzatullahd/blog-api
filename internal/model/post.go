package model

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
