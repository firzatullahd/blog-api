package model

type Post struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status"`
}

type FilterFindPost struct {
	IDs    []uint64 `db:"id"`
	Limit  int
	Offset int
}

type FilterFindRPost struct {
	PostID uint64   `db:"post_id"`
	TagIDs []uint64 `db:"tag_id"`
}

type FilterSearchPost struct {
	TagLabel []string
}
