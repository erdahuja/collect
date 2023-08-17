package question

import (
	"collect/foundation/event"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	ErrNotFound              = errors.New("question not found")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of APIs for user access.
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

func (c *Core) Create(ctx context.Context, np NewQuestion) (Question, error) {
	now := time.Now()
	que := Question{
		FormID:       np.FormID,
		QuestionType: np.QuestionType,
		QuestionText: np.QuestionText,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tran := func(s Storer) error {
		quesDB, err := s.Create(ctx, que)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		que = quesDB
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return Question{}, fmt.Errorf("tran: %w", err)
	}

	if err := c.evnCore.SendEvent(np.CreatedEvent(que)); err != nil {
		return Question{}, fmt.Errorf("failed to send a `%s` event: %w", EventCreated, err)
	}


	return que, nil
}

func (c *Core) QueryByFormID(ctx context.Context, formID int64) ([]Question, error) {
	ques, err := c.storer.QueryByFormID(ctx, formID)
	if err != nil {
		return []Question{}, fmt.Errorf("query: %w", err)
	}

	return ques, nil
}

func (c *Core) QueryByID(ctx context.Context, quesID int64) (Question, error) {
	que, err := c.storer.QueryByID(ctx, quesID)
	if err != nil {
		return Question{}, fmt.Errorf("query: %w", err)
	}

	return que, nil
}

func (c *Core) Delete(ctx context.Context, quesID int64) error {
	err := c.storer.Delete(ctx, quesID)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}
