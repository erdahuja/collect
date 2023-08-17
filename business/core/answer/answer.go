package answer

import (
	"collect/business/core/response"
	"collect/foundation/event"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("response not found")

// Core manages the set of APIs for user access.
type Core struct {
	log     *zap.SugaredLogger
	storer  Storer
	evnCore *event.Core
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, evnCore *event.Core, storer Storer) *Core {
	c := Core{
		log:     log,
		storer:  storer,
		evnCore: evnCore,
	}
	c.registerEventHandlers(evnCore)
	return &c
}

// Create adds a Answer to the database.
func (c *Core) Create(ctx context.Context, np NewAnswer, res response.Response) (Answer, error) {
	now := time.Now()
	ans := Answer{
		QuestionID: np.QuestionID,
		ResponseID: res.ResponseID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if np.AnswerOptionID >= 0 {
		ans.AnswerOptionID = &np.AnswerOptionID
	}
	if np.AnswerText != "" {
		ans.AnswerText = &np.AnswerText
	}

	tran := func(s Storer) error {
		ad, err := s.Create(ctx, ans)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		ans.AnswerID = ad.AnswerID
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return Answer{}, fmt.Errorf("tran: %w", err)
	}

	if err := c.evnCore.SendEvent(np.CreatedEvent(ans, res)); err != nil {
		return Answer{}, fmt.Errorf("failed to send a `%s` event: %w", EventCreated, err)
	}

	return ans, nil
}

// QueryByFormID gets all the responses of a form from the database.
func (c *Core) QueryByFormID(ctx context.Context, formID int64) ([]Answer, error) {
	responses, err := c.storer.QueryByFormID(ctx, formID)
	if err != nil {
		return []Answer{}, fmt.Errorf("query: %w", err)
	}

	return responses, nil
}

// QueryByFormAndRespondentID gets the response for a user to a form
func (c *Core) QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (Answer, error) {
	responses, err := c.storer.QueryByFormAndRespondentID(ctx, formID, respID)
	if err != nil {
		return Answer{}, fmt.Errorf("query: %w", err)
	}

	return responses, nil
}
