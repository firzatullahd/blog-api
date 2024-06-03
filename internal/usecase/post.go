package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	customerror "github.com/firzatullahd/blog-api/internal/model/error"

	"github.com/firzatullahd/blog-api/internal/entity"
	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/utils"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func (u *Usecase) CreatePost(ctx context.Context, in *model.Post) (*entity.Post, error) {
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
			tx.Rollback()
		}
	}()

	if len(missingTags) > 0 {
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

	post.Tags = in.Tags

	return post, tx.Commit()
}

func (u *Usecase) UpdatePost(ctx context.Context, in *model.Post, id uint64, email string) error {
	logCtx := fmt.Sprintf("%T.UpdatePost", u)

	// if publish post, check user's role
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

	rposts, err := u.repo.FindRPostTag(ctx, model.FilterFindRPost{PostID: id})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	var existingTagIDs []uint64
	for _, v := range rposts {
		existingTagIDs = append(existingTagIDs, v.TagID)
	}

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{
		ID: existingTagIDs,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return err
	}

	existingTags := make([]string, len(tags))
	mapTagsID := make(map[string]uint64)
	for _, v := range tags {
		existingTags = append(existingTags, v.Label)
		mapTagsID[v.Label] = v.ID
	}

	fmt.Printf("MAPTAG %v\n", mapTagsID)

	missingTags := utils.FindMissing(existingTags, in.Tags)
	unusedTags := utils.FindMissing(in.Tags, existingTags)

	missingTagIDs := make([]uint64, 0)
	if len(missingTags) > 0 {
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
			tagId, ok := mapTagsID[v]
			if ok {
				err = u.repo.DeleteRPostTag(ctx, tx, id, tagId)
				if err != nil {
					logger.Error(ctx, logCtx, err)
					return err
				}
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

func (u *Usecase) GetPost(ctx context.Context, id uint64) (*entity.Post, error) {
	logCtx := fmt.Sprintf("%T.GetPost", u)

	posts, err := u.repo.FindPosts(ctx, model.FilterFindPost{
		IDs: []uint64{id},
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customerror.ErrNotFound
		}
		return nil, err
	}

	if len(posts) == 0 {
		return nil, customerror.ErrNotFound
	}
	post := posts[0]

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

	return &post, nil
}

func (u *Usecase) SearchPost(ctx context.Context, in model.FilterSearchPost) ([]entity.Post, int, error) {
	logCtx := fmt.Sprintf("%T.SearchPost", u)
	// API Search Post
	// find tags based on labels inputted, get tag id
	// if labels inputted empty, then get all tags.
	// find post ids from r_post_tag based on tag id
	// find posts, map with each tags.
	// add pagination

	tags, err := u.repo.FindTag(ctx, model.FilterFindTag{Label: in.TagLabel})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, 0, err
	}

	mapTags := make(map[uint64]string)
	var tagIDs []uint64
	for _, v := range tags {
		tagIDs = append(tagIDs, v.ID)
		mapTags[v.ID] = v.Label
	}

	rposts, err := u.repo.FindRPostTag(ctx, model.FilterFindRPost{TagIDs: tagIDs})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, 0, err
	}

	mapPosts := make(map[uint64][]string)
	var postIDs []uint64
	for _, v := range rposts {
		postIDs = append(postIDs, v.PostID)
		mapPosts[v.PostID] = append(mapPosts[v.PostID], mapTags[v.TagID])
	}

	posts, err := u.repo.FindPosts(ctx, model.FilterFindPost{
		IDs: postIDs,
	})
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, 0, err
	}

	countPosts, err := u.repo.CountPost(ctx)
	if err != nil {
		logger.Error(ctx, logCtx, err)
		return nil, 0, err
	}

	for i := range posts {
		posts[i].Tags = mapPosts[posts[i].ID]
	}

	return posts, countPosts, nil
}
