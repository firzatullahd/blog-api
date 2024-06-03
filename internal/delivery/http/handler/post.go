package handler

import (
	"encoding/json"
	"fmt"
	"github.com/firzatullahd/blog-api/internal/delivery/http/middleware"
	"github.com/firzatullahd/blog-api/internal/model"
	customerror "github.com/firzatullahd/blog-api/internal/model/error"
	"github.com/firzatullahd/blog-api/internal/model/response"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.Post

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.SetHTTPResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	post, err := h.Usecase.CreatePost(ctx, &payload)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusCreated, "success create post", post)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	logCtx := fmt.Sprintf("%T.UpdatePost", h)
	ctx := r.Context()
	var payload model.Post
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error(ctx, logCtx, err)
		response.SetHTTPResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	data, ok := r.Context().Value(middleware.UserDataKey).(model.UserData)
	if !ok {
		logger.Error(ctx, logCtx, fmt.Errorf("error"))
		response.SetHTTPResponse(w, http.StatusUnauthorized, customerror.ErrUnauthorized.Error(), nil)
		return
	}

	id := r.PathValue("id")
	postId, err := strconv.ParseUint(id, 10, 64)

	err = h.Usecase.UpdatePost(ctx, &payload, postId, data.Email)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		logger.Error(ctx, logCtx, err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusOK, "success update post", nil)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	logCtx := fmt.Sprintf("%T.DeletePost", h)
	ctx := r.Context()

	id := r.PathValue("id")
	postId, err := strconv.ParseUint(id, 10, 64)

	err = h.Usecase.DeletePost(ctx, postId)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		logger.Error(ctx, logCtx, err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusOK, "success delete post", nil)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	logCtx := fmt.Sprintf("%T.GetPost", h)
	ctx := r.Context()

	id := r.PathValue("id")
	postId, err := strconv.ParseUint(id, 10, 64)

	post, err := h.Usecase.GetPost(ctx, postId)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		logger.Error(ctx, logCtx, err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusOK, "record found", post)
}

func (h *Handler) SearchPost(w http.ResponseWriter, r *http.Request) {
	logCtx := fmt.Sprintf("%T.SearchPost", h)
	ctx := r.Context()

	tag := r.URL.Query().Get("tag")
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")

	var payload model.FilterSearchPost
	payload.Limit, _ = strconv.Atoi(limit)
	payload.Page, _ = strconv.Atoi(page)
	if len(tag) > 0 {
		payload.TagLabel = strings.Split(tag, ",")
	}

	if payload.Page == 0 {
		payload.Page = 1
	}

	post, err := h.Usecase.SearchPost(ctx, payload)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		logger.Error(ctx, logCtx, err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusOK, "record found", post)
}
