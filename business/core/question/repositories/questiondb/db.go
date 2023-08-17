package questiondb

import (
	"bytes"
	"collect/business/core/question"
	"collect/business/sys/database"
	"context"
	"fmt"
	"strings"

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

// WithinTran runs passed function and do commit/rollback at the end.
func (s *Store) WithinTran(ctx context.Context, fn func(s question.Storer) error) error {
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

// Create adds a Question to the database. It returns the created Form with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, que question.Question) (question.Question, error) {
	const q = `
	INSERT INTO questions
		(form_id, question_type, question_text, created_at, updated_at)
	VALUES
	    (:form_id, :question_type, :question_text, :created_at, :updated_at)
		`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBQuestion(que)); err != nil {
		return question.Question{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	var id int64
	qs := `select nextval('questions_question_id_seq'); `
	if err := database.QueryRowContext(ctx, s.log, s.db, qs, &id); err != nil {
		return question.Question{}, fmt.Errorf("QueryRowContext: %v", err)
	}

	que.ID = id - 1

	return que, nil
}

// QueryByFormID gets the specified user from the database.
func (s *Store) QueryByFormID(ctx context.Context, formID int64) ([]question.Question, error) {
	data := struct {
		ID int64 `db:"form_id"`
	}{
		ID: formID,
	}
	const q = `
	SELECT
		*
	FROM
		questions
		where form_id = :form_id
	`

	var dbQuestions []dbQues
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbQuestions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreQuestionSlice(dbQuestions), nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, quesID int64) (question.Question, error) {

	data := struct {
		ID int64 `db:"question_id"`
	}{
		ID: quesID,
	}

	const q = `
	SELECT
		*
	FROM
		questions
		WHERE
		question_id = :question_id
	`

	var dbQuestions dbQues
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbQuestions); err != nil {
		return question.Question{}, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreQuestion(dbQuestions), nil
}

// Delete removes the product identified by a given ID.
func (s *Store) Delete(ctx context.Context, id int64) error {
	data := struct {
		ID int64 `db:"question_id"`
	}{
		ID: id,
	}

	const q = `
	DELETE FROM
		questions
	WHERE
		question_id = :question_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) applyFilter(filter question.QueryFilter, buf *bytes.Buffer) {
	var wc []string

	if filter.FormID != nil {
		wc = append(wc, " f.form_id = :form_id")
	}

	if filter.ID != nil {
		wc = append(wc, " f.question_id = :question_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
