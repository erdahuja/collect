package usercache

import (
	"context"
	"strconv"
	"sync"

	"collect/business/core/user"

	"go.uber.org/zap"
)

// =============================================================================

// Store manages the set of APIs for user data and caching.
type Store struct {
	log    *zap.SugaredLogger
	storer user.Storer
	cache  map[string]*user.User
	mu     sync.RWMutex
}

// NewStore constructs the api for data and caching access.
func NewStore(log *zap.SugaredLogger, storer user.Storer) *Store {
	return &Store{
		log:    log,
		storer: storer,
		cache:  map[string]*user.User{},
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s *Store) WithinTran(ctx context.Context, fn func(s user.Storer) error) error {
	return s.storer.WithinTran(ctx, fn)
}

// Create inserts a new user into the database.
func (s *Store) Create(ctx context.Context, usr user.User) (user.User, error) {
	ud, err := s.storer.Create(ctx, usr)
	if err != nil {
		return usr, err
	}

	s.writeCache(usr)

	return ud, nil
}

// Update replaces a user document in the database.
func (s *Store) Update(ctx context.Context, usr user.User) error {
	if err := s.storer.Update(ctx, usr); err != nil {
		return err
	}

	s.writeCache(usr)

	return nil
}

// Delete removes a user from the database.
func (s *Store) Delete(ctx context.Context, usr user.User) error {
	if err := s.storer.Delete(ctx, usr); err != nil {
		return err
	}

	s.deleteCache(usr)

	return nil
}

// Query retrieves a list of existing users from the database.
func (s *Store) Query(ctx context.Context) ([]user.User, error) {
	return s.storer.Query(ctx)
}

// QueryByID gets the specified user from the database.
func (s *Store) QueryByID(ctx context.Context, userID int) (user.User, error) {
	cachedUsr, ok := s.readCache(strconv.FormatInt(int64(userID), 10))
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.storer.QueryByID(ctx, userID)
	if err != nil {
		return user.User{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// QueryByEmail gets the specified user from the database by email.
func (s *Store) QueryByEmail(ctx context.Context, email string) (user.User, error) {
	cachedUsr, ok := s.readCache(email)
	if ok {
		return cachedUsr, nil
	}

	usr, err := s.storer.QueryByEmail(ctx, email)
	if err != nil {
		return user.User{}, err
	}

	s.writeCache(usr)

	return usr, nil
}

// =============================================================================

// readCache performs a safe search in the cache for the specified key.
func (s *Store) readCache(key string) (user.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	usr, exists := s.cache[key]
	if !exists {
		return user.User{}, false
	}

	return *usr, true
}

// writeCache performs a safe write to the cache for the specified user.
func (s *Store) writeCache(usr user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[strconv.FormatInt(usr.ID, 10)] = &usr
	s.cache[usr.Email] = &usr
}

// deleteCache performs a safe removal from the cache for the specified user.
func (s *Store) deleteCache(usr user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, strconv.FormatInt(usr.ID, 10))
	delete(s.cache, usr.Email)
}
