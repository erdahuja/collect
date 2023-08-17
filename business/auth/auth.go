package auth

import (
	"collect/business/core/user"
	"collect/business/core/user/repositories/userdb"
	"collect/foundation/event"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// These the current set of rules we have for auth.
const (
	RuleAuthenticate  = "auth"
	RuleAny           = "allowAny"
	RuleAdminOnly     = "allowOnlyAdmin"
	RuleUserOnly      = "allowOnlyUser"
	RuleCollectorOnly = "allowOnlyCollector"
)

// ErrForbidden is returned when a auth issue is identified.
var ErrForbidden = errors.New("attempted action is not allowed")

// Config represents information required to initialize auth.
type Config struct {
	Log    *zap.SugaredLogger
	DB     *sqlx.DB
	Secret []byte
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log    *zap.SugaredLogger
	user   *user.Core
	method jwt.SigningMethod
	parser *jwt.Parser
	secret []byte
}

// New creates an Auth to support authentication/authorization.
func New(event *event.Core, cfg Config) *Auth {
	// If a database connection is not provided, we won't perform the
	// user active check.
	var usr *user.Core
	if cfg.DB != nil {
		usr = user.NewCore(event, userdb.NewStore(cfg.Log, cfg.DB))
	}

	a := Auth{
		log:    cfg.Log,
		user:   usr,
		secret: cfg.Secret,
		method: jwt.GetSigningMethod(jwt.SigningMethodHS256.Name),
		parser: jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name})),
	}

	return &a
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		a.log.Infof("roles: %v, subject: %s", claims.Roles, claims.Subject)
		// Check the database for this user to verify they are still active.
		if !a.isUserEnabled(ctx, *claims) {
			return Claims{}, fmt.Errorf("user not active : %w", err)
		}
		return *claims, nil
	} else {
		fmt.Println(err)
		return *claims, err
	}

}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(ctx context.Context, claims Claims, rule string) error {

	if rule == RuleAdminOnly {
		if ok := a.IsAdmin(claims.Roles); ok {
			return nil
		}
	} else if rule == RuleUserOnly {
		if ok := a.IsUser(claims.Roles); ok {
			return nil
		}
	} else if rule == RuleAny {
		if ok := a.IsAny(claims.Roles); ok {
			return nil
		}
	} else if rule == RuleCollectorOnly {
		if ok := a.IsCollector(claims.Roles); ok {
			return nil
		}
	}

	return ErrForbidden
}

func (a *Auth) IsAdmin(roles []user.Role) bool {
	for _, role := range roles {
		if role == user.RoleAdmin {
			return true
		}
	}
	return false
}

func (a *Auth) IsUser(roles []user.Role) bool {
	for _, role := range roles {
		if role == user.RoleUser {
			return true
		}
	}
	return false
}

func (a *Auth) IsCollector(roles []user.Role) bool {
	for _, role := range roles {
		if role == user.RoleCollector {
			return true
		}
	}
	return false
}

func (a *Auth) IsAny(roles []user.Role) bool {
	for _, role := range roles {
		if role == user.RoleAdmin || role == user.RoleUser || role == user.RoleCollector {
			return true
		}
	}
	return false
}

// isUserEnabled hits the database and checks the user is not disabled. If the
// no database connection was provided, this check is skipped.
func (a *Auth) isUserEnabled(ctx context.Context, claims Claims) bool {
	if a.user == nil {
		return true
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return false
	}

	usr, err := a.user.QueryByID(ctx, userID)
	if err != nil {
		return false
	}

	return usr.Active
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)

	str, err := token.SignedString(a.secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}
