package responsedb

import (
	"collect/business/core/response"
	"collect/business/sys/database"
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user database access.
type Store struct {
	log    *zap.SugaredLogger
	db     sqlx.ExtContext
	inTran bool
}

// NewStore constructs the api for data access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create adds a Response to the database. It returns the created response with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, resp response.Response) (response.Response, error) {
	const q = `
	INSERT INTO responses
		(form_id, respondent_id, created_at, updated_at)
	VALUES
		(:form_id, :respondent_id, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBResponse(resp)); err != nil {
		return response.Response{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	var id int64
	qs := `select nextval('responses_response_id_seq'); `
	if err := database.QueryRowContext(ctx, s.log, s.db, qs, &id); err != nil {
		return response.Response{}, fmt.Errorf("QueryRowContext: %v", err)
	}

	resp.ResponseID = id - 1

	return resp, nil
}

// QueryByFormAndRespondentID find in db if respondent has already responded to the form
func (s *Store) QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (response.Response, error) {
	data := struct {
		FormID       int64  `db:"form_id"`
		RespondentID string `db:"respondent_id"`
	}{
		FormID:       formID,
		RespondentID: respID,
	}
	const q = `
	Select * 
	from
	responses
		WHERE 
		(form_id = :form_id AND respondent_id = :respondent_id)`

	var resp dbResponse
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &resp); err != nil {
		return response.Response{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	return toCoreResponse(resp), nil
}

// QueryByID find in db if respondent has already responded to the form
func (s *Store) QueryByID(ctx context.Context, respID int64) (response.Response, error) {
	data := struct {
		ResponseID int64 `db:"response_id"`
	}{
		ResponseID: respID,
	}
	const q = `
	Select * 
	from
	responses
		WHERE 
		(response_id = :response_id)`

	var resp dbResponse
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &resp); err != nil {
		return response.Response{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	return toCoreResponse(resp), nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByFormID(ctx context.Context, formID int64) ([]response.Response, error) {
	data := struct {
		FormID int64 `db:"form_id"`
	}{
		FormID: formID,
	}

	const q = `
	SELECT
		*
	FROM
		responses
	WHERE
		form_id = :form_id`

	var ci []dbResponse
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &ci); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []response.Response{}, err
		}
		return []response.Response{}, fmt.Errorf("selecting formID[%q]: %w", formID, err)
	}

	return toCoreResponseSlice(ci), nil
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s *Store) WithinTran(ctx context.Context, fn func(s response.Storer) error) error {
	if s.inTran {
		return fn(s)
	}

	f := func(tx *sqlx.Tx) error {
		s := &Store{
			log:    s.log,
			db:     tx,
			inTran: true,
		}
		return fn(s)
	}

	return database.WithinTran(ctx, s.log, s.db.(*sqlx.DB), f)
}
