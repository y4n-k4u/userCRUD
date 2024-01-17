package persistence

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
	"userCRUD/internal/common"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/model"
	"userCRUD/pkg/common/password"
)

var (
	ErrUserNotFound      = errors.New("user with the ID not found")
	ErrUserAlreadyExists = errors.New("user already exists with given ID")
	ErrUsernameTaken     = errors.New("username is already taken")
	ErrEmailTaken        = errors.New("email is already taken")
	ErrPageOutOfRange    = errors.New("page out of range")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UpdateUser) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, pagination *common.Pagination) ([]*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, rawPassword string) (*model.User, error)
}

type UserRepositoryMemory struct {
	sync.RWMutex
	l              deps.Logger
	orderedUserIDs []string
	usersByID      map[string]*model.User
	userByUsername map[string]*model.User
	userByEmail    map[string]*model.User
}

func NewUserRepositoryMemory(l deps.Logger) *UserRepositoryMemory {

	ur := &UserRepositoryMemory{
		l:              l,
		orderedUserIDs: make([]string, 0, 16),
		usersByID:      make(map[string]*model.User),
		userByUsername: make(map[string]*model.User),
		userByEmail:    make(map[string]*model.User),
	}

	adminPassHashed, _ := password.HashPassword("admin")

	ur.CreateUser(context.Background(), &model.User{
		ID:       uuid.New().String(),
		Email:    "admin@gmail.com",
		Username: "admin",
		Password: adminPassHashed,
		Admin:    true,
	})

	return ur
}

func (r *UserRepositoryMemory) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.usersByID[user.ID]; exists {
		return nil, ErrUserAlreadyExists
	}

	if _, emailExists := r.userByEmail[user.Email]; emailExists {
		return nil, ErrEmailTaken
	}

	if _, usernameExists := r.userByUsername[user.Username]; usernameExists {
		return nil, ErrUsernameTaken
	}

	r.usersByID[user.ID] = user
	r.userByUsername[user.Username] = user
	r.userByEmail[user.Email] = user
	r.orderedUserIDs = append(r.orderedUserIDs, user.ID)

	r.l.Info(ctx, "User created", "id", user.ID)

	return user, nil
}

func (r *UserRepositoryMemory) UpdateUser(ctx context.Context, userU *model.UpdateUser) (*model.User, error) {
	r.Lock()
	defer r.Unlock()

	exUser, found := r.usersByID[userU.ID]
	if !found {
		return nil, ErrUserNotFound
	}

	if userU.Email != exUser.Email && userU.Email != "" {
		if _, emailExists := r.userByEmail[userU.Email]; emailExists {
			return nil, ErrEmailTaken
		}
	}

	if userU.Username != exUser.Username && userU.Username != "" {
		if _, usernameExists := r.userByUsername[userU.Username]; usernameExists {
			return nil, ErrUsernameTaken
		}
	}

	user := &model.User{
		ID:       userU.ID,
		Email:    userU.Email,
		Username: userU.Username,
		Password: userU.Password,
		Admin:    userU.Admin,
	}

	delete(r.userByUsername, exUser.Username)
	delete(r.userByEmail, exUser.Email)

	exUser = user
	r.userByUsername[user.Username] = user
	r.userByEmail[user.Email] = user

	r.l.Info(ctx, "User updated", "id", user.ID)

	return user, nil
}

func (r *UserRepositoryMemory) DeleteUser(ctx context.Context, id string) error {
	r.Lock()
	defer r.Unlock()

	exUser, found := r.usersByID[id]
	if !found {
		return ErrUserNotFound
	}

	for i, userID := range r.orderedUserIDs {
		if userID == id {
			r.orderedUserIDs = append(r.orderedUserIDs[:i], r.orderedUserIDs[i+1:]...)
			break
		}
	}

	delete(r.userByUsername, exUser.Username)
	delete(r.userByEmail, exUser.Email)
	delete(r.usersByID, id)

	r.l.Info(ctx, "User deleted", "id", id)

	return nil
}

func (r *UserRepositoryMemory) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.usersByID[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepositoryMemory) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.userByUsername[username]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepositoryMemory) GetUsers(ctx context.Context, pagination *common.Pagination) ([]*model.User, error) {
	r.RLock()
	defer r.RUnlock()

	start := int((pagination.Page - 1) * pagination.PageSize)
	end := start + int(pagination.PageSize)

	if start >= len(r.orderedUserIDs) {
		return nil, ErrPageOutOfRange
	}

	if end > len(r.orderedUserIDs) {
		end = len(r.orderedUserIDs)
	}

	users := make([]*model.User, 0, end-start)
	for _, userID := range r.orderedUserIDs[start:end] {
		users = append(users, r.usersByID[userID])
	}

	return users, nil
}

func (r *UserRepositoryMemory) GetUserByUsernameAndPassword(ctx context.Context, username, rawPassword string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, err := r.GetUserByUsername(ctx, username)
	if user == nil && err != nil {
		return nil, ErrUserNotFound
	}

	if !password.CheckPasswordHash(rawPassword, user.Password) {
		return nil, ErrUserNotFound
	}

	return user, nil
}
