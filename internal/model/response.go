package model

type Response[T any] struct {
	Data    T             `json:"data"`
	Message string        `json:"message"`
	Paging  *PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type GeneralResponse struct {
	Message string `json:"message"`
}
