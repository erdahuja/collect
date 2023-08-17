package user

import (
	"errors"
	"regexp"
)

// Used to check for sql injection problems.
var sqlInjection = regexp.MustCompile("^[A-Za-z0-9_]+$")

// =============================================================================

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID    *int64  `validate:"omitempty,uuid4"`
	Name  *string `validate:"omitempty,min=3"`
	Email *string `validate:"omitempty,email"`
}

// ByID sets the ID field of the QueryFilter value.
func (f *QueryFilter) ByID(id int64) {
	var zero int64
	if id != zero {
		f.ID = &id
	}
}

// ByName sets the Name field of the QueryFilter value.
func (f *QueryFilter) ByName(name string) error {
	if name != "" {
		if !sqlInjection.MatchString(name) {
			return errors.New("invalid name format")
		}

		f.Name = &name
	}

	return nil
}

// ByEmail sets the Email field of the QueryFilter value.
func (f *QueryFilter) ByEmail(email string) {
	var zero string
	if email != zero {
		f.Email = &email
	}
}
