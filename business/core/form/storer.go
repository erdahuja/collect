package form

import "context"

type Storer interface {
	WithinTran(ctx context.Context, fn func(s Storer) error) error
	Create(ctx context.Context, itm Form) (Form, error)
	Query(ctx context.Context, filter QueryFilter) ([]Form, error)
	QueryByID(ctx context.Context, formID int64) (Form, error)
	Update(ctx context.Context, itm Form) error
	Delete(ctx context.Context, filter QueryFilter) error
}
