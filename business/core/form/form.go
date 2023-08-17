// package form contains form related CRUD functionality.
package form

import (
	"collect/business/sys/validate"
	"collect/foundation/event"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("user not found")

// Core manages the set of APIs for form access.
type Core struct {
	log     *zap.SugaredLogger
	storer  Storer
	evnCore *event.Core
}

// NewCore constructs a core for form api access.
func NewCore(log *zap.SugaredLogger, evnCore *event.Core, storer Storer) *Core {
	core := Core{
		log:     log,
		storer:  storer,
		evnCore: evnCore,
	}
	core.registerEventHandlers(evnCore)
	return &core
}

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	FormID *int `validate:"omitempty,int"`
}

// ByID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) ByID(q int) {
	var zero int
	if q != zero {
		f.FormID = &q
	}
}

// WithID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) WithID(id int) {
	f.FormID = &id
}

// Create adds a Form to the database. It returns the created Form with
// fields like ID and DateCreated populated.
func (c *Core) Create(ctx context.Context, np NewForm) (Form, error) {
	now := time.Now()
	frm := Form{
		Title:       np.Title,
		Description: np.Description,
		DateCreated: now,
		DateUpdated: now,
	}

	tran := func(s Storer) error {
		fm, err := s.Create(ctx, frm)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		frm.ID = fm.ID
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return Form{}, fmt.Errorf("tran: %w", err)
	}

	if err := c.evnCore.SendEvent(np.CreatedEvent(frm)); err != nil {
		return Form{}, fmt.Errorf("failed to send a `%s` event: %w", EventCreated, err)
	}

	return frm, nil
}

// Query gets all Products from the database.
func (c *Core) Query(ctx context.Context, filter QueryFilter) ([]Form, error) {
	prds, err := c.storer.Query(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// QueryByID gets the specified form from the database.
func (c *Core) QueryByID(ctx context.Context, formID int64) (Form, error) {
	form, err := c.storer.QueryByID(ctx, formID)
	if err != nil {
		return Form{}, fmt.Errorf("query: %w", err)
	}

	return form, nil
}

// Delete removed the specified form from the database along with all responses.
func (c *Core) Delete(ctx context.Context, filter QueryFilter) error {
	err := c.storer.Delete(ctx, filter)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

// Update replaces a user document in the database.
func (c *Core) Update(ctx context.Context, itm Form, uu UpdateForm) (Form, error) {
	if err := validate.Check(uu); err != nil {
		return Form{}, fmt.Errorf("validating data: %w", err)
	}

	if uu.Title != nil {
		itm.Title = *uu.Title
	}
	if uu.Description != nil {
		itm.Description = *uu.Description
	}
	itm.DateUpdated = time.Now()

	tran := func(s Storer) error {
		if err := s.Update(ctx, itm); err != nil {
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return Form{}, fmt.Errorf("tran: %w", err)
	}

	return itm, nil
}
