package persistence

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/model"
)

var ctx = context.TODO()
var logger = &deps.MockLogger{}

func TestCreateAndUpdateUser(t *testing.T) {
	ur := NewUserRepositoryMemory(&deps.MockLogger{})

	userFirst := &model.User{ID: uuid.New().String(), Username: "userFirst", Email: "userFirst@example.com"}
	userSecond := &model.User{ID: uuid.New().String(), Username: "userSecond", Email: "userSecond@example.com"}
	user, err := ur.CreateUser(ctx, userFirst)
	if err != nil || user == nil {
		t.Errorf("Failed to add user: %v", err)
	}

	user, err = ur.CreateUser(ctx, userSecond)
	if err != nil || user == nil {
		t.Errorf("Failed to create user: %v", err)
	}

	user, err = ur.CreateUser(ctx, userFirst)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}

	userSecondUpdate := &model.UpdateUser{ID: userSecond.ID, Username: userFirst.Username, Email: "newtest@example.com"}
	user, err = ur.UpdateUser(ctx, userSecondUpdate)
	if !errors.Is(err, ErrUsernameTaken) {
		t.Errorf("Expected ErrUsernameTaken, got %v", err)
	}

	userSecondUpdate = &model.UpdateUser{ID: userSecond.ID, Username: "newuser", Email: "userFirst@example.com"}
	user, err = ur.UpdateUser(ctx, userSecondUpdate)
	if !errors.Is(err, ErrEmailTaken) {
		t.Errorf("Expected ErrEmailTaken, got %v", err)
	}

	userSecondUpdate = &model.UpdateUser{ID: userSecond.ID, Username: "newuser", Email: "newtest@example.com"}
	user, err = ur.UpdateUser(ctx, userSecondUpdate)
	if err != nil || user == nil {
		t.Errorf("Failed to update user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	ur := NewUserRepositoryMemory(logger)

	userFirst := &model.User{ID: uuid.New().String(), Username: "userFirst", Email: "userFirst@example.com"}
	user, err := ur.CreateUser(ctx, userFirst)
	if err != nil || user == nil {
		t.Errorf("Failed to create user: %v", err)
	}

	err = ur.DeleteUser(ctx, uuid.New().String())
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}

	err = ur.DeleteUser(ctx, userFirst.ID)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}
}
