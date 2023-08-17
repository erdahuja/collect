package formgrp

import (
	"collect/business/auth"
	"collect/business/core/form"
	"collect/business/core/question"
	"collect/business/core/user"
	"collect/business/sys/validate"
	"collect/foundation/web"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of form endpoints.
type Handlers struct {
	Form     *form.Core
	Question *question.Core
	Auth     *auth.Auth
}

// Create adds a new user to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu form.NewForm
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.Form.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return web.NewRequestError(err, http.StatusConflict)
		}
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

// Query returns a list of forms
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter, err := parseFilter(r)
	if err != nil {
		return err
	}
	users, err := h.Form.Query(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to query for users: %w", err)
	}
	return web.Respond(ctx, w, users, http.StatusOK)
}

// Delete removed of form
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter, err := parseFilter(r)
	if err != nil {
		return err
	}
	err = h.Form.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to query for users: %w", err)
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}

func parseFilter(r *http.Request) (form.QueryFilter, error) {
	values := r.URL.Query()

	var filter form.QueryFilter

	if formID := values.Get("form_id"); formID != "" {
		qua, err := strconv.ParseInt(formID, 10, 64)
		if err != nil {
			return form.QueryFilter{}, validate.NewFieldsError("formID", err)
		}
		filter.WithID(int(qua))
	}

	return filter, nil
}

// QueryQuestionsByFormID returns a user by its ID.
func (h Handlers) QueryQuestionsByFormID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	formID := web.Param(r, "form_id")
	if formID == "" {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	ques, err := h.Question.QueryByFormID(ctx, id)
	if err != nil {
		if !errors.Is(err, user.ErrNotFound) {
			return fmt.Errorf("ID[%v]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, ques, http.StatusOK)
}
