package core_domain

import "time"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponce struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ID         string    `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func NewRegisterResonce(user User) *RegisterResponce {
	return &RegisterResponce{
		Name:       user.Name,
		Email:      user.Email,
		ID:         user.ID,
		Created_at: user.Created_at,
		Updated_at: user.Updated_at,
	}
}

type User struct {
	Name       string
	Email      string
	Password   string
	ID         string
	Created_at time.Time
	Updated_at time.Time
}

// func NewUser(
// 	name string,
// 	email string,
// 	password string,
// 	id int,
// 	version int,
// ) User {
// 	return User{
// 		Name:     name,
// 		Email:    email,
// 		Password: password,
// 		ID:       id,
// 		Version:  version,
// 	}
// }
