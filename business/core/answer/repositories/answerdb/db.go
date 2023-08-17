package answerdb

import (
	"collect/business/core/answer"
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
func (s *Store) Create(ctx context.Context, ans answer.Answer) (answer.Answer, error) {
	const q = `
	INSERT INTO answers
		(question_id, response_id, answer_text, answer_option_id, created_at, updated_at)
	VALUES
		(:question_id, :response_id, :answer_text, :answer_option_id, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBAnswer(ans)); err != nil {
		return answer.Answer{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	var id int64
	qs := `select nextval('answers_answer_id_seq'); `
	if err := database.QueryRowContext(ctx, s.log, s.db, qs, &id); err != nil {
		return answer.Answer{}, fmt.Errorf("QueryRowContext: %v", err)
	}

	ans.AnswerID = id

	return ans, nil
}

// QueryByFormAndRespondentID find in db if respondent has already responded to the form
func (s *Store) QueryByFormAndRespondentID(ctx context.Context, formID int64, respID string) (answer.Answer, error) {
	data := struct {
		FormID       int64  `db:"form_id"`
		RespondentID string `db:"respondent_id"`
	}{
		FormID:       formID,
		RespondentID: respID,
	}
	const q = `
	Select response_id
	from
	answers
		WHERE 
		(form_id = :form_id AND respondent_id = :respondent_id)`

	var resp answer.Answer
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &resp); err != nil {
		return answer.Answer{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	return resp, nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByFormID(ctx context.Context, formID int64) ([]answer.Answer, error) {
	data := struct {
		FormID int64 `db:"form_id"`
	}{
		FormID: formID,
	}

	const q = `
	SELECT
		*
	FROM
		answers
	WHERE
		form_id = :form_id`

	var ci []dbAnswer
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &ci); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []answer.Answer{}, err
		}
		return []answer.Answer{}, fmt.Errorf("selecting formID[%q]: %w", formID, err)
	}

	return toCoreAnswerSlice(ci), nil
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s *Store) WithinTran(ctx context.Context, fn func(s answer.Storer) error) error {
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
