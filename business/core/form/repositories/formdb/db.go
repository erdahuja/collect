package formdb

import (
	"bytes"
	"collect/business/core/form"
	"collect/business/sys/database"
	"context"
	"errors"
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

// Create adds a Form to the database.
func (s *Store) Create(ctx context.Context, frm form.Form) (form.Form, error) {
	const q = `
	INSERT INTO forms
		(form_title, form_description, created_at, updated_at)
	VALUES
		(:form_title, :form_description, :created_at, :updated_at)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBForm(frm)); err != nil {
		return form.Form{}, fmt.Errorf("namedexeccontext: %w", err)
	}

	var id int64
	qs := `select nextval('forms_form_id_seq'); `
	if err := database.QueryRowContext(ctx, s.log, s.db, qs, &id); err != nil {
		return form.Form{}, fmt.Errorf("QueryRowContext: %v", err)
	}

	frm.ID = id - 1

	return frm, nil
}

// Query gets all Forms from the database.
func (s *Store) Query(ctx context.Context, filter form.QueryFilter) ([]form.Form, error) {

	const q = `SELECT * FROM forms f`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, buf)

	var dbForms []dbForm
	if err := database.NamedQuerySlice(ctx, s.log, s.db, buf.String(), struct{}{}, &dbForms); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreFormSlice(dbForms), nil
}

// Update fixes all forms in the database.
func (s *Store) Update(ctx context.Context, itm form.Form) error {
	const q = `
	UPDATE
		forms
	SET
		"title" = :form_title,
		"form_description" = :form_description
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBForm(itm)); err != nil {
		return fmt.Errorf("updating userID[%d]: %w", itm.ID, err)
	}

	return nil
}

func (s *Store) applyFilter(filter form.QueryFilter, buf *bytes.Buffer) {
	var wc []string

	if filter.FormID != nil {
		wc = append(wc, " f.form_id = :form_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, formID int64) (form.Form, error) {
	data := struct {
		FormID int64 `db:"form_id"`
	}{
		FormID: formID,
	}

	const q = `
	SELECT
		*
	FROM
		forms
	WHERE
		form_id = :form_id`

	var usr dbForm
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return form.Form{}, form.ErrNotFound
		}
		return form.Form{}, fmt.Errorf("selecting formID[%q]: %w", formID, err)
	}

	return toCoreForm(usr), nil
}

// Delete removed the specified form from the database.
func (s *Store) Delete(ctx context.Context, filter form.QueryFilter) error {

	const q = `
	delete
	FROM
		forms`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, buf)

	if err := database.NamedExecContext(ctx, s.log, s.db, buf.String(), struct{}{}); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return form.ErrNotFound
		}
		return fmt.Errorf("selecting formID[%v]: %w", filter.FormID, err)
	}

	return nil
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s *Store) WithinTran(ctx context.Context, fn func(s form.Storer) error) error {
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
