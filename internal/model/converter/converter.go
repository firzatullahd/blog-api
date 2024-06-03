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
