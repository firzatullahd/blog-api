package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/model/converter"
	customerror "github.com/firzatullahd/blog-api/internal/model/error"
	"time"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (u *Usecase) CreatePost(ctx context.Context, in *model.Post) (*model.PostResult, error) {
	logCtx := fmt.Sprintf("%T.CreatePost", u)

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{
		Label: in.Tags,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	tagIDs := make([]uint64, 0)
	originalTags := make([]string, len(tags))
	for _, v := range tags {
		originalTags = append(originalTags, v.Label)
		tagIDs = append(tagIDs, v.ID)
	}

	missingTags := utils.FindMissing(originalTags, in.Tags)

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// if all tags exists, continue
	// if any tag not exists, insert missing tag
	if len(missingTags) > 0 {
		// insert tags
		for _, v := range missingTags {
			tagId, err := u.repo.InsertTag(ctx, tx, &entity.Tag{
				Label: v,
			})
			if err != nil {
				logger.Error(ctx, logCtx, err)
				return nil, err
			}

			tagIDs = append(tagIDs, tagId)
		}
	}

	post, err := u.repo.InsertPost(ctx, tx, &entity.Post{
		Title:   in.Title,
		Content: in.Content,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	for _, v := range tagIDs {
		err := u.repo.InsertRPostTag(ctx, tx, post.ID, v)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return nil, err
		}
	}

	_ = tx.Commit()

	post.Tags = in.Tags
	resp := converter.EntityPostToResponse(post)
	return &resp, nil
}

func (u *Usecase) UpdatePost(ctx context.Context, in *model.Post, id uint64, email string) error {
	logCtx := fmt.Sprintf("%T.UpdatePost", u)

	if len(in.Status) > 0 {
		users, err := u.repo.FindUsers(ctx, &model.FilterFindUser{
			Email: &email,
		})
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return err
		}

		if len(users) > 0 {
			return customerror.ErrNotFound
		}

		if users[0].Role != entity.RoleAdmin.String() {
			return customerror.ErrForbidden
		}
	}

	//rposts, err := u.repo.FindRPostTag(ctx, model.FilterFindRPost{PostID: []uint64{id}})
	//if err != nil {
	//	logger.Error(ctx, logCtx, err)
	//	return err
	//}

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{
		Label: in.Tags,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	originalTags := make([]string, len(tags))
	mapTags := make(map[string]uint64)
	for _, v := range tags {
		originalTags = append(originalTags, v.Label)
		mapTags[v.Label] = v.ID
	}

	fmt.Printf("MAPTAG %v\n", mapTags)

	missingTags := utils.FindMissing(originalTags, in.Tags)
	unusedTags := utils.FindMissing(in.Tags, originalTags)

	tx, err := u.repo.WithTransaction()
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	updatePayload := entity.Post{
		ID:      id,
		Title:   in.Title,
		Content: in.Content,
	}
	switch in.Status {
	case entity.StatusPublish.String():
		updatePayload.Status = entity.StatusPublish.String()
		now := time.Now()
		updatePayload.PublishDate = &now
	default:
		updatePayload.Status = entity.StatusDraft.String()
		updatePayload.PublishDate = nil
	}

	err = u.repo.UpdatePost(ctx, tx, &updatePayload)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	missingTagIDs := make([]uint64, 0)
	if len(missingTags) > 0 {
		// insert tags
		for _, v := range missingTags {
			tagId, err := u.repo.InsertTag(ctx, tx, &entity.Tag{
				Label: v,
			})
			if err != nil {
				logger.Error(ctx, logCtx, err)
				return err
			}

			missingTagIDs = append(missingTagIDs, tagId)
		}
	}

	for _, v := range missingTagIDs {
		err := u.repo.InsertRPostTag(ctx, tx, id, v)
		if err != nil {
			logger.Error(ctx, logCtx, err)
			return err
		}
	}

	if len(unusedTags) > 0 {
		for _, v := range unusedTags {
			tagId, _ := mapTags[v]
			err = u.repo.DeleteRPostTag(ctx, tx, id, tagId)
			if err != nil {
				logger.Error(ctx, logCtx, err)
				return err
			}
		}
	}

	return tx.Commit()
}

func (u *Usecase) DeletePost(ctx context.Context, id uint64) error {
	logCtx := fmt.Sprintf("%T.DeletePost", u)
	tx, err := u.repo.WithTransaction()

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = u.repo.DeletePost(ctx, tx, id)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	return tx.Commit()
}

func (u *Usecase) GetPost(ctx context.Context, id uint64) (*model.PostResult, error) {
	logCtx := fmt.Sprintf("%T.GetPost", u)
	post, err := u.repo.FindPost(ctx, model.FilterFindPost{
		ID: id,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}

	rposts, err := u.repo.FindRPostTag(ctx, model.FilterFindRPost{PostID: id})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	var tagIDs []uint64
	for _, v := range rposts {
		tagIDs = append(tagIDs, v.TagID)
	}

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{ID: tagIDs})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	for _, tag := range tags {
		post.Tags = append(post.Tags, tag.Label)
	}
	post.TagIDs = tagIDs
	resp := converter.EntityPostToResponse(post)
	return &resp, nil
}

func (u *Usecase) SearchPost(ctx context.Context, tagLabels []string) ([]model.PostResult, error) {
	logCtx := fmt.Sprintf("%T.SearchPost", u)

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{
		Label: tagLabels,
	})

	var tagIDs []uint64
	for _, v := range tags {
		tagIDs = append(tagIDs, v.ID)
	}

	// find rposts
	// find post

	post, err := u.repo.FindPost(ctx, model.FilterFindPost{
		ID: id,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}

	rposts, err := u.repo.FindRPostTag(ctx, model.FilterFindRPost{PostID: id})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	var tagIDs []uint64
	for _, v := range rposts {
		tagIDs = append(tagIDs, v.TagID)
	}

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{ID: tagIDs})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, err
	}

	for _, tag := range tags {
		post.Tags = append(post.Tags, tag.Label)
	}
	post.TagIDs = tagIDs
	resp := converter.EntityPostToResponse(post)
	return &resp, nil
}
