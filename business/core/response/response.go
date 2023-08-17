package response

import (
	"collect/business/sys/database"
	"collect/foundation/event"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("response not found")

// Core manages the set of APIs for response access.
type Core struct {
	log     *zap.SugaredLogger
	storer  Storer
	evnCore *event.Core
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, evnCore *event.Core, storer Storer) *Core {
	core := Core{
		log:     log,
		storer:  storer,
		evnCore: evnCore,
	}
	core.registerEventHandlers(evnCore)
	return &core
}

// Create adds a Response to the database.
func (c *Core) Create(ctx context.Context, np NewResponse) (Response, error) {
	now := time.Now()
	res := Response{
		FormID:       np.FormID,
		RespondentID: np.RespondentID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tran := func(s Storer) error {
		respd, err := s.Create(ctx, res)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		res.ResponseID = respd.ResponseID
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return Response{}, fmt.Errorf("tran: %w", err)
	}

	if err := c.evnCore.SendEvent(np.CreatedEvent(res)); err != nil {
		return Response{}, fmt.Errorf("failed to send a `%s` event: %w", EventCreated, err)
	}

	return res, nil
}

// QueryByFormID gets all the responses of a form from the database.
func (c *Core) QueryByFormID(ctx context.Context, formID int64) ([]Response, error) {
	responses, err := c.storer.QueryByFormID(ctx, formID)
	if err != nil {
		return []Response{}, fmt.Errorf("query: %w", err)
	}

	return responses, nil
}

// QueryByID gets the response from the database.
func (c *Core) QueryByID(ctx context.Context, respID int64) (Response, error) {
	responses, err := c.storer.QueryByID(ctx, respID)
	if err != nil {
		if !errors.Is(err, database.ErrDBNotFound) {
			return Response{}, fmt.Errorf("query: %w", err)
		}
	}

	return responses, nil
}

// QueryByFormAndRespondentID gets the response for a user to a form
func (c *Core) QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (Response, error) {
	response, err := c.storer.QueryByFormAndRespondentID(ctx, formID, respID)
	if err != nil {
		if !errors.Is(err, database.ErrDBNotFound) {
			return Response{}, fmt.Errorf("query: %w", err)
		}
	}

	return response, nil
}
