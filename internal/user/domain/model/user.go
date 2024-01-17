package model

type User struct {
	ID       string `validate:"uuid4"`
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=5"`
	Password string `validate:"required,min=5"`
	Admin    bool   `validate:"boolean"`
}
type UpdateUser struct {
	ID       string `validate:"required,uuid4"`
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=5"`
	Password string `validate:"required,min=5"`
	Admin    bool   `validate:"required,boolean"`
}

type UserByID struct {
	ID string `validate:"uuid4"`
}

type UserByUsername struct {
	Username string `validate:"required,min=5"`
}

type Users struct {
	Page int `validate:"required,min=1"`
	Size int `validate:"required,min=1,max=50"`
}
