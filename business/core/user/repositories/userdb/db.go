package userdb

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"collect/business/core/user"

	database "collect/business/sys/database"

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
func (s *Store) WithinTran(ctx context.Context, fn func(s user.Storer) error) error {
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

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) (user.User, error) {
	const q = `
	INSERT INTO users
		( name, email, roles, password_hash, active, date_created, date_updated)
	VALUES
		(:name, :email, :roles, :password_hash, :active, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return usr, fmt.Errorf("create: %w", user.ErrUniqueEmail)
		}
		return usr, fmt.Errorf("inserting user: %w", err)
	}

	var id int64
	qs := `select nextval('users_user_id_seq'); `
	if err := database.QueryRowContext(ctx, s.log, s.db, qs, &id); err != nil {
		return user.User{}, fmt.Errorf("QueryRowContext: %v", err)
	}

	usr.ID = id - 1
	return usr, nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr user.User) error {
	const q = `
	UPDATE
		users
	SET
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"date_updated" = :date_updated,
		active = :active
	WHERE
		id = :id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, database.ErrDBDuplicatedEntry) {
			return user.ErrUniqueEmail
		}
		return fmt.Errorf("updating userID[%d]: %w", usr.ID, err)
	}

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: strconv.FormatInt(usr.ID, 10),
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting userID[%d]: %w", usr.ID, err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context) ([]user.User, error) {

	const q = `
	SELECT
		*
	FROM
		users
	`

	var usrs []dbUser
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, struct{}{}, &usrs); err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return toCoreUserSlice(usrs), nil
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, id int) (user.User, error) {
	data := struct {
		ID int `db:"user_id"`
	}{
		ID: id,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		user_id = :user_id`

	var usr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("selecting userID[%q]: %w", id, err)
	}

	return toCoreUser(usr), nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = :email`

	var usr dbUser
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usr); err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}

	return toCoreUser(usr), nil
}
