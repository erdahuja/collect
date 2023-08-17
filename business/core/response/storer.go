package response

import "context"

type Storer interface {
	WithinTran(ctx context.Context, fn func(s Storer) error) error
	Create(ctx context.Context, res Response) (Response, error)
	QueryByID(ctx context.Context, respID int64) (Response, error)
	QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (Response, error)
	QueryByFormID(ctx context.Context, formID int64) ([]Response, error)
}
