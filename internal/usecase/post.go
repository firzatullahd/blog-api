package usecase

import (
	"context"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (u *Usecase) CreatePost(ctx context.Context, in *model.Post) error {
	logCtx := fmt.Sprintf("%T.CreatePost", u)

	// find tags
	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{
		Label: in.Tags,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	tagIDs := make([]uint64, 0)
	originalTags := make([]string, len(tags))
	for _, v := range tags {
		originalTags = append(originalTags, v.Label)
		tagIDs = append(tagIDs, v.ID)
	}

	_, missingTags := utils.ArrayCompare(originalTags, in.Tags)

	tx, err := u.repo.WithTransaction()
	if err != nil {
		return err
	}
	// if all tags exists, continue
	// if any tag not exists, insert missing tag
	if len(missingTags) > 0 {
		// insert tags
		for _, v := range missingTags {
			tagId, err := u.repo.InsertTag(ctx, tx, &entity.Tag{
				Label: v,
			})
			if err != nil {
				return err
			}

			tagIDs = append(tagIDs, tagId)
		}
	}

	postId, err := u.repo.InsertPost(ctx, tx, &entity.Post{
		Title:   in.Title,
		Content: in.Content,
	})
	if err != nil {
		return err
	}

	// update tags with post id
	for _, v := range tagIDs {
		err := u.repo.InsertRPostTag(ctx, tx, postId, v)
		if err != nil {
			return err
		}
	}

	return nil
}
