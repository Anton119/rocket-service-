package iam_test

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/iam/internal/errors"
	"github.com/Anton119/rocket-service-/iam/internal/model"
)

type userRepoFake struct {
	mu       sync.Mutex
	byLogin  map[string]model.User
	byUUID   map[uuid.UUID]model.User
	createFn func(ctx context.Context, user model.User) error
}

func newUserRepoFake() *userRepoFake {
	return &userRepoFake{
		byLogin: make(map[string]model.User),
		byUUID:  make(map[uuid.UUID]model.User),
	}
}

func (f *userRepoFake) Create(ctx context.Context, user model.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.createFn != nil {
		return f.createFn(ctx, user)
	}

	if _, ok := f.byLogin[user.Login]; ok {
		return errs.ErrUserAlreadyExists
	}

	f.byLogin[user.Login] = user
	f.byUUID[user.UUID] = user

	return nil
}

func (f *userRepoFake) GetByUUID(_ context.Context, id uuid.UUID) (model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	user, ok := f.byUUID[id]
	if !ok {
		return model.User{}, errs.ErrUserNotFound
	}

	return user, nil
}

func (f *userRepoFake) GetByLogin(_ context.Context, login string) (model.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	user, ok := f.byLogin[login]
	if !ok {
		return model.User{}, errs.ErrUserNotFound
	}

	return user, nil
}

type sessionRepoFake struct {
	mu       sync.Mutex
	sessions map[uuid.UUID]model.Session
	createFn func(ctx context.Context, session model.Session, ttl time.Duration) error
	getFn    func(ctx context.Context, sessionUUID uuid.UUID) (model.Session, error)
	deleteFn func(ctx context.Context, sessionUUID uuid.UUID) error
}

func newSessionRepoFake() *sessionRepoFake {
	return &sessionRepoFake{
		sessions: make(map[uuid.UUID]model.Session),
	}
}

func (f *sessionRepoFake) Create(_ context.Context, session model.Session, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.createFn != nil {
		return f.createFn(context.Background(), session, 0)
	}

	f.sessions[session.UUID] = session

	return nil
}

func (f *sessionRepoFake) Get(_ context.Context, sessionUUID uuid.UUID) (model.Session, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.getFn != nil {
		return f.getFn(context.Background(), sessionUUID)
	}

	session, ok := f.sessions[sessionUUID]
	if !ok {
		return model.Session{}, errs.ErrSessionNotFound
	}

	return session, nil
}

func (f *sessionRepoFake) Delete(_ context.Context, sessionUUID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.deleteFn != nil {
		return f.deleteFn(context.Background(), sessionUUID)
	}

	delete(f.sessions, sessionUUID)

	return nil
}
