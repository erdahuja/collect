package user

import "errors"

// Set of possible roles for a user.
var (
	RoleAdmin     = Role{"ADMIN"}
	RoleUser      = Role{"USER"}
	RoleCollector = Role{"COLLECTOR"}
)

// Set of known roles.
var roles = map[string]Role{
	RoleAdmin.name:     RoleAdmin,
	RoleUser.name:      RoleUser,
	RoleCollector.name: RoleCollector,
}

// Role represents a role in the system.
type Role struct {
	name string
}

// ParseRole parses the string value and returns a role if one exists.
func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, errors.New("invalid role")
	}

	return role, nil
}

// Name returns the name of the role.
func (r Role) Name() string {
	return r.name
}

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (r *Role) UnmarshalText(data []byte) error {
	r.name = string(data)
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

// Equal provides support for the go-cmp package and testing.
func (r Role) Equal(r2 Role) bool {
	return r.name == r2.name
}
