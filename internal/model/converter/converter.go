package converter

import (
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
)

func EntityPostToResponse(in *entity.Post) model.PostResult {
	return model.PostResult{
		ID:          in.ID,
		Title:       in.Title,
		Content:     in.Content,
		Status:      in.Status,
		PublishDate: in.PublishDate,
		Tags:        in.Tags,
	}
}

func SearchPostResponse(posts []entity.Post, count int, page int) model.SearchResult {
	var data []model.PostResult
	for _, v := range posts {
		data = append(data, EntityPostToResponse(&v))
	}
	return model.SearchResult{
		Data:  data,
		Count: count,
		Page:  page,
	}
}
