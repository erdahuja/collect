package answer

import "context"

type Storer interface {
	WithinTran(ctx context.Context, fn func(s Storer) error) error
	Create(ctx context.Context, res Answer) (Answer, error)
	QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (Answer, error)
	QueryByFormID(ctx context.Context, formID int64) ([]Answer, error)
}
