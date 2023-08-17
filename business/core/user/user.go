// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"collect/business/sys/validate"
	"collect/foundation/event"

	"golang.org/x/crypto/bcrypt"
)

// =============================================================================

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidEmail          = errors.New("email is not valid")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of APIs for user access.
type Core struct {
	storer  Storer
	evnCore *event.Core
}

// NewCore constructs a core for user api access.
func NewCore(evnCore *event.Core, storer Storer) *Core {

	core := Core{
		storer:  storer,
		evnCore: evnCore,
	}
	return &core
}

// Create inserts a new user into the database.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	now := time.Now()

	usr := User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		Active:       true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	tran := func(s Storer) error {
		ud, err := s.Create(ctx, usr)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		usr.ID = ud.ID
		return nil
	}

	if err := c.storer.WithinTran(ctx, tran); err != nil {
		return User{}, fmt.Errorf("tran: %w", err)
	}

	if err := c.evnCore.SendEvent(nu.CreatedEvent(usr)); err != nil {
		return User{}, fmt.Errorf("failed to send a `%s` event: %w", EventCreated, err)
	}

	return usr, nil
}

// Query retrieves a list of existing users from the database.
func (c *Core) Query(ctx context.Context) ([]User, error) {

	users, err := c.storer.Query(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// QueryByID gets the specified user from the database.
func (c *Core) QueryByID(ctx context.Context, userID int) (User, error) {
	user, err := c.storer.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	if err := c.evnCore.SendEvent(user.QueryEvent()); err != nil {
		return User{}, fmt.Errorf("failed to send a `%s` event: %w", EventQuery, err)
	}

	return user, nil
}

// QueryByEmail gets the specified user from the database by email.
func (c *Core) QueryByEmail(ctx context.Context, email string) (User, error) {
	user, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	return user, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c *Core) Authenticate(ctx context.Context, email string, password string) (User, error) {
	usr, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, ErrAuthenticationFailure
	}

	return usr, nil
}
