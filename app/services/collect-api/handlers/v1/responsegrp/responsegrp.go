package responsegrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"collect/business/core/answer"
	"collect/business/core/response"
	"collect/foundation/web"
)

// Set of error variables for handling user group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	Response *response.Core
	Answer   *answer.Core
}

// Create adds a new response to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu response.NewResponse
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	resp, err := h.Response.QueryByFormAndRespondentID(ctx, nu.FormID, nu.RespondentID)
	if err != nil {
		return fmt.Errorf("response[%+v]: %w", &nu, err)
	}
	if resp.ResponseID == 0 {
		resp, err = h.Response.Create(ctx, nu)
		if err != nil {
			return fmt.Errorf("response[%+v]: %w", &resp, err)
		}
	}
	return web.Respond(ctx, w, resp, http.StatusCreated)
}

// QueryByFormID returns responses for a form.
func (h Handlers) QueryByFormID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	formID := web.Param(r, "form_id")
	if formID == "" {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(formID, 10, 64)
	if err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	responses, err := h.Response.QueryByFormID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", formID, err)
		}
	}

	return web.Respond(ctx, w, responses, http.StatusOK)
}

// CreateAnswer adds a new response to the system.
func (h Handlers) CreateAnswer(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	responseID := web.Param(r, "id")
	if responseID == "" {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(responseID, 10, 64)
	if err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	resp, err := h.Response.QueryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("response[%+d]: %w", id, err)
	}

	if resp.ResponseID == 0 {
		return web.Respond(ctx, w, fmt.Errorf("response id is invalid[%+v]: %w", &id, err), http.StatusNotFound)
	}

	var nu answer.NewAnswer
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	ans, err := h.Answer.Create(ctx, nu, resp)
	if err != nil {
		return fmt.Errorf("response[%+v]: %w", &resp, err)
	}

	return web.Respond(ctx, w, ans, http.StatusCreated)
}
