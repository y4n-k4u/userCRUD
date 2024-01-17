package command

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"userCRUD/internal/common"
	"userCRUD/internal/common/constants"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/model"
	"userCRUD/internal/user/infrastructure/persistence"
	"userCRUD/pkg/common/password"
)

var (
	ErrAuthFailed           = errors.New("authentication failed")
	ErrNotEnoughPermissions = errors.New("operation requires more privileges")
)

type User struct {
	ur        persistence.UserRepository
	validator deps.Validator
}

func NewUserCommand(ur persistence.UserRepository, v deps.Validator) *User {
	return &User{
		ur:        ur,
		validator: v,
	}
}

func (u *User) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if !isAdmin(ctx) {
		return nil, ErrNotEnoughPermissions
	}

	user.ID = uuid.New().String()

	if err := u.validator.Struct(user); err != nil {
		return nil, err
	}

	hashedPass, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPass

	if user, err = u.ur.CreateUser(ctx, user); err != nil || user == nil {
		return nil, err
	}

	return user, nil
}

func (u *User) UpdateUser(ctx context.Context, userU *model.UpdateUser) (*model.User, error) {
	if !isAdmin(ctx) {
		return nil, ErrNotEnoughPermissions
	}

	if err := u.validator.Struct(userU); err != nil {
		return nil, err
	}

	if userU.Password != "" {
		hashedPass, err := password.HashPassword(userU.Password)
		if err != nil {
			return nil, err
		}
		userU.Password = hashedPass
	}

	user, err := u.ur.UpdateUser(ctx, userU)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) DeleteUser(ctx context.Context, userID *model.UserByID) error {
	if !isAdmin(ctx) {
		return ErrNotEnoughPermissions
	}

	if err := u.validator.Struct(userID); err != nil {
		return err
	}

	if err := u.ur.DeleteUser(ctx, userID.ID); err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByID(ctx context.Context, userID *model.UserByID) (*model.User, error) {
	if err := u.validator.Struct(userID); err != nil {
		return nil, err
	}

	user, err := u.ur.GetUserByID(ctx, userID.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUserByUsername(ctx context.Context, username *model.UserByUsername) (*model.User, error) {
	if err := u.validator.Struct(username); err != nil {
		return nil, err
	}

	user, err := u.ur.GetUserByUsername(ctx, username.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUsers(ctx context.Context, pagination *common.Pagination) ([]*model.User, error) {
	if err := u.validator.Struct(pagination); err != nil {
		return nil, err
	}

	users, err := u.ur.GetUsers(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func isAdmin(ctx context.Context) bool {
	ctxUser, ok := ctx.Value(constants.UserContextKey).(*model.User)
	if !ok {
		return false
	}

	if !ctxUser.Admin {
		return false
	}

	return true
}
