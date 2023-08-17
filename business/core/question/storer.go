package question

import "context"

type Storer interface {
	WithinTran(ctx context.Context, fn func(s Storer) error) error
	Create(ctx context.Context, que Question) (Question, error)
	QueryByFormID(ctx context.Context, formID int64) ([]Question, error)
	QueryByID(ctx context.Context, quesID int64) (Question, error)
	Delete(ctx context.Context, quesID int64) error
}

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	FormID *int `validate:"omitempty,int"`
	ID     *int `validate:"omitempty,int"`
}

// ByFormID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) ByFormID(q int) {
	var zero int
	if q != zero {
		f.FormID = &q
	}
}

// WithFormID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) WithFormID(id int) {
	f.FormID = &id
}

// ByID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) ByID(q int) {
	var zero int
	if q != zero {
		f.ID = &q
	}
}

// WithID sets the FormID field of the QueryFilter value.
func (f *QueryFilter) WithID(id int) {
	f.ID = &id
}
