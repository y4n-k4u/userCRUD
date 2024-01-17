package command

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"testing"
	"userCRUD/internal/common/constants"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/model"
	"userCRUD/internal/user/infrastructure/persistence"
)

var (
	ur    = persistence.NewUserRepositoryMemory(&deps.MockLogger{})
	admin = &model.User{
		Email:    "admin@gmail.com",
		Username: "admin",
		Password: "admin",
		Admin:    true,
	}
	command = NewUserCommand(ur, deps.NewGoPlaygroundValidator())
)

func TestCreateUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), constants.UserContextKey, admin)

	newUser := model.User{
		Username: "newUser",
		Email:    "newUser@gmail.com",
		Password: "password",
		Admin:    false,
	}

	user, err := command.CreateUser(ctx, &newUser)

	if err != nil {
		t.Errorf("Failed to create new user: %s", err)
	}
	if user.ID == "" {
		t.Errorf("User ID is empty")
	}
	if user.Password == "password" {
		t.Errorf("User password isn't hashed")
	}

	invalidUser := model.User{
		Username: "User",
		Email:    "newUser.gmail.com",
		Password: "pass",
		Admin:    false,
	}

	user, err = command.CreateUser(ctx, &invalidUser)

	if user != nil {
		t.Errorf("The user must be nil")
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		if len(validationErrs) != 3 {
			t.Errorf("Expected 3 validation errors, but got %d", len(validationErrs))
		}
	} else {
		t.Errorf("Expected validation error, got different type of error")
	}
}
