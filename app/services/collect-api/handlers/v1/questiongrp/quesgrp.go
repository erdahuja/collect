// question api system
package questiongrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"collect/business/core/question"
	"collect/foundation/web"
)

// Set of error variables for handling question group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of questions endpoints.
type Handlers struct {
	Question *question.Core
}

// Create adds a new question to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu question.NewQuestion
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	crt, err := h.Question.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("question[%+v]: %w", &crt, err)
	}

	return web.Respond(ctx, w, crt, http.StatusCreated)
}

// QueryByID returns a question by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	questionID := web.Param(r, "question_id")
	if questionID == "" {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(questionID, 10, 64)
	if err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	que, err := h.Question.QueryByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, question.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", questionID, err)
		}
	}

	return web.Respond(ctx, w, que, http.StatusOK)
}

// DeleteQuestion deleted a question.
func (h Handlers) DeleteQuestion(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu question.DeleteQuestion
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	que, err := h.Question.QueryByID(ctx, nu.ID)
	if err != nil {
		return web.NewRequestError(errors.New("unable to find question:"+err.Error()), http.StatusNotFound)
	}

	err = h.Question.Delete(ctx, que.ID)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &que, err)
	}

	return web.Respond(ctx, w, que, http.StatusOK)
}
