package handler

import (
	"encoding/json"
	"net/http"

	"github.com/firzatullahd/blog-api/internal/model"
	customerror "github.com/firzatullahd/blog-api/internal/model/error"
	"github.com/firzatullahd/blog-api/internal/model/response"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.SetHTTPResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	resp, err := h.Usecase.Register(ctx, &payload)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusCreated, "User registered successfully", resp)
	return
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.SetHTTPResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	resp, err := h.Usecase.Login(ctx, &payload)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}
	response.SetHTTPResponse(w, http.StatusOK, "User logged in successfully", resp)
	return
}

func (h *Handler) GrantAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.SetHTTPResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	resp, err := h.Usecase.GrantAdmin(ctx, &payload)
	if err != nil {
		code, errMsg := customerror.ParseError(err)
		response.SetHTTPResponse(w, code, errMsg, nil)
		return
	}

	response.SetHTTPResponse(w, http.StatusOK, "successfully granted admin role", resp)
	return
}
